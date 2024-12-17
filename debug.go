package echogy

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/echogy-io/echogy/pkg/logger"
	"github.com/echogy-io/echogy/pkg/stat"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io/fs"
	"net"
	"net/http"
	"sync"
)

type fakeListen struct {
	chs        chan gossh.Channel
	serverConn *gossh.ServerConn
}

func (f *fakeListen) Close() error {
	close(f.chs)
	return nil
}

func (f *fakeListen) Accept() (net.Conn, error) {
	ch := <-f.chs
	if nil != ch {
		return wrapChannelConn(f.serverConn, ch), nil
	}
	return nil, errors.New("error ")
}

func (f *fakeListen) Addr() net.Addr {
	return f.serverConn.LocalAddr()
}

func newListen(conn *gossh.ServerConn) *fakeListen {
	return &fakeListen{
		chs:        make(chan gossh.Channel, 2),
		serverConn: conn,
	}
}

//go:embed debugger/dist/*
var webContent embed.FS

type debugServer struct {
	ctx      ssh.Context
	fake     *fakeListen
	chEvent  chan *EventMessage
	server   *http.Server
	isClosed bool
}

func (f *debugServer) Close() error {
	if !f.isClosed {
		f.isClosed = true
		f.server.Close()
	}
	return nil
}

type EventMessage struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Request struct {
	Method  string            `json:"method"`
	Uri     string            `json:"uri"`
	Headers map[string]string `json:"headers"`
}

type Response struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
}

type WrapHttpEntity struct {
	Id       int       `json:"id"`
	Response *Response `json:"response"`
	Request  *Request  `json:"request"`
	UseTime  int64     `json:"useTime"`
}

type Stats struct {
	RequestBytes      int64 `json:"requestBytes"`
	ResponseBytes     int64 `json:"responseBytes"`
	Requests          int   `json:"requests"`
	Responses         int   `json:"responses"`
	ActiveConnections int   `json:"activeConnections"`
	TotalConnections  int   `json:"totalConnections"`
}

type SyncEventMessage struct {
	Tunnel       string            `json:"tunnel"`
	Stats        *Stats            `json:"stats"`
	HttpEntities []*WrapHttpEntity `json:"httpEntities"`
}

type UpdateEventMessage struct {
	Stats      *Stats          `json:"stats"`
	HttpEntity *WrapHttpEntity `json:"httpEntity"`
}

func writeEvent(w http.ResponseWriter, msg *EventMessage) error {
	if fl, ok := w.(http.Flusher); ok {
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		w.Write([]byte("event: sync\n"))
		fmt.Fprintf(w, "data: %s\n\n", data)
		fl.Flush()
	}
	return nil
}

func (f *debugServer) getStat() *Stats {
	s := stat.GetStat(f.ctx)
	return &Stats{
		RequestBytes:      s.Receive,
		ResponseBytes:     s.Send,
		Requests:          s.Request,
		Responses:         s.Response,
		ActiveConnections: s.ConnCount,
		TotalConnections:  s.TotalConn,
	}
}

func (f *debugServer) eventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queue := stat.GetQueue(f.ctx)

	sem := &SyncEventMessage{
		Stats:        f.getStat(),
		HttpEntities: make([]*WrapHttpEntity, 0),
	}

	if t, ok := f.ctx.Value(sshTunnelAddrKey).(string); ok {
		sem.Tunnel = t
	}

	for i, item := range queue.Items() {
		e := item.(*stat.RequestEntity)
		entity := &WrapHttpEntity{
			Id:       i + 1,
			Request:  simpleRequest(e.Request),
			Response: simpleResponse(e.Response),
			UseTime:  e.UseTime,
		}
		sem.HttpEntities = append(sem.HttpEntities, entity)
	}

	err := writeEvent(w, &EventMessage{
		Name: "all",
		Data: sem,
	})

	if nil != err {
		logger.Error("send event [all] error", err, map[string]interface{}{
			"server": "debug",
		})
	}

	for {
		select {
		case <-r.Context().Done():
			return
		case <-f.ctx.Done():
			return
		case msg := <-f.chEvent:
			if nil != msg {
				err := writeEvent(w, msg)
				logger.Error(fmt.Sprintf("send event [%s] error", msg.Name), err, map[string]interface{}{
					"server": "debug",
				})
			}
		}
	}
}

func (f *debugServer) newChan(ch gossh.Channel) {
	if f.isClosed {
		logger.WarnN("debug server is closed")
		return
	}
	if nil != ch {
		f.fake.chs <- ch
	}
}

func simpleHeader(h http.Header) map[string]string {
	sh := make(map[string]string)
	for k, v := range h {
		sh[k] = v[0]
	}
	return sh
}

func simpleRequest(r *http.Request) *Request {
	h := simpleHeader(r.Header)
	h["Host"] = r.Host
	return &Request{
		Method:  r.Method,
		Uri:     r.RequestURI,
		Headers: h,
	}
}

func simpleResponse(w *http.Response) *Response {
	return &Response{
		Status:  w.StatusCode,
		Headers: simpleHeader(w.Header),
	}
}

func (f *debugServer) UpdateEvent(w *http.Response, r *http.Request, t int64) {
	m := &EventMessage{
		Name: "update",
		Data: &UpdateEventMessage{
			Stats: f.getStat(),
			HttpEntity: &WrapHttpEntity{
				Request:  simpleRequest(r),
				Response: simpleResponse(w),
				UseTime:  t,
			},
		},
	}
	f.chEvent <- m
}

func getOrCreateDebugServer(conn *gossh.ServerConn, ctx ssh.Context) *debugServer {
	var svr *debugServer
	var ok bool
	if svr, ok = ctx.Value(sshDebugServer).(*debugServer); !ok {
		ln := newListen(conn)
		svr = &debugServer{
			fake:    ln,
			ctx:     ctx,
			chEvent: make(chan *EventMessage, 2),
		}

		wg := sync.WaitGroup{}

		dist, err := fs.Sub(webContent, "debugger/dist")
		if err != nil {
			return nil
		}
		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.FS(dist)))
		mux.Handle("/events", http.HandlerFunc(svr.eventHandler))

		server := &http.Server{
			Handler: mux,
		}
		svr.server = server

		ctx.SetValue(sshDebugServer, svr)
		wg.Add(1)
		go func() {
			wg.Done()
			err = server.Serve(svr.fake)
			if nil != err {
				logger.Error("debug server error", err, map[string]interface{}{
					"server": "debug",
				})
			}
		}()
		wg.Wait()
	}
	return svr
}
