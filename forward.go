package echogy

import (
	"context"
	"fmt"
	"github.com/echogy-io/echogy/pkg/logger"
	"github.com/echogy-io/echogy/pkg/stat"
	"github.com/echogy-io/echogy/pkg/tui"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type fwdConn struct {
	ch   gossh.Channel
	conn net.Conn
}

func (f *fwdConn) Close() error {
	if nil != f.ch {
		f.ch.Close()
	}
	if nil != f.conn {
		f.conn.Close()
	}
	return nil
}

type forwarder struct {
	context           context.Context
	cancelFunc        context.CancelFunc
	sess              ssh.Session
	accessId          string
	pty               *tui.HttpReversProxyPty
	remoteForwardChan chan net.Conn
	chanCounter       atomic.Int64
	chanMap           *sync.Map
}

func newForwarder(accessId, domain string, session ssh.Session) (*forwarder, error) {
	pty, err := tui.NewHttpReverseProxyPty(session, fmt.Sprintf("%s.%s", accessId, domain))
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithCancel(session.Context())
	return &forwarder{
		context:           ctx,
		cancelFunc:        cancelFunc,
		accessId:          accessId,
		pty:               pty,
		sess:              session,
		chanMap:           &sync.Map{},
		remoteForwardChan: make(chan net.Conn, 4),
	}, nil
}

type remoteForwardChannelData struct {
	DestAddr   string
	DestPort   uint32
	OriginAddr string
	OriginPort uint32
}

func (fwd *forwarder) dispatchRemoteForward(hijackConn *hijackHttp) {
	hijackConn.SetDispatch(func(w *http.Response, r *http.Request, t int64) {
		stat.Put(fwd.sess.Context(), w, r, t)
		fwd.pty.Update()
		if debug, ok := fwd.sess.Context().Value(sshDebugServer).(*debugServer); ok {
			debug.UpdateEvent(w, r, t)
		}
	})
	fwd.remoteForwardChan <- hijackConn
}

func (fwd *forwarder) serve() {
	remoteAddr := fwd.sess.RemoteAddr().String()
	logger.Info("created dispatchRemoteForward conn", map[string]interface{}{
		"module":     "conn",
		"accessId":   fwd.accessId,
		"remoteAddr": remoteAddr,
	})

	go func() {
		err := fwd.pty.Start()
		if err != nil {
			logger.Error("start pty conn", err, map[string]interface{}{
				"module":     "conn",
				"accessId":   fwd.accessId,
				"remoteAddr": remoteAddr,
			})
		}
		fwd.cancelFunc()
	}()

	keepalive := time.NewTicker(30 * time.Second)

	defer keepalive.Stop()

	counter := 0

	for {
		select {
		case <-fwd.context.Done():
			return
		case <-keepalive.C:
			if counter > 5 {
				fwd.cancelFunc()
			} else {
				_, err := fwd.sess.SendRequest("keepalive@openssh.com", true, nil)
				if err != nil {
					logger.WarnN("Failed to send keepalive request")
					counter++
				} else {
					logger.DebugN("send keepalive request")
					counter = 0
				}
			}
		case facadeConn := <-fwd.remoteForwardChan:
			go fwd.doRemoteForwarded(facadeConn)
		}
	}
}

func (fwd *forwarder) getForwardDest() *remoteForwardRequest {
	_, localPortStr, _ := net.SplitHostPort(fwd.sess.LocalAddr().String())
	localPort, _ := strconv.Atoi(localPortStr)
	value := fwd.sess.Context().Value(sshRemoteForward)
	if nil != value {
		return value.(*remoteForwardRequest)
	}
	return &remoteForwardRequest{
		BindAddr: "localhost",
		BindPort: uint32(localPort),
	}
}

func (fwd *forwarder) doRemoteForwarded(facadeConn net.Conn) {
	s := stat.GetStat(fwd.sess.Context())
	s.ConnCount += 1
	s.TotalConn += 1
	defer func() {
		s.ConnCount -= 1
	}()

	remoteAddr := fwd.sess.RemoteAddr().String()
	svrConn := fwd.sess.Context().Value(ssh.ContextKeyConn).(*gossh.ServerConn)
	logger.Debug("open dispatchRemoteForward channel", map[string]interface{}{
		"module":     "conn",
		"accessId":   fwd.accessId,
		"remoteAddr": remoteAddr,
	})

	facadeRequestAddr, facadeRequestPortStr, _ := net.SplitHostPort(facadeConn.RemoteAddr().String())
	facadePort, _ := strconv.Atoi(facadeRequestPortStr)
	dest := fwd.getForwardDest()

	payload := gossh.Marshal(&remoteForwardChannelData{
		DestAddr:   dest.BindAddr,
		DestPort:   dest.BindPort,
		OriginAddr: facadeRequestAddr,
		OriginPort: uint32(facadePort),
	})

	gosshChan, _, err := svrConn.OpenChannel("forwarded-tcpip", payload)

	if err != nil {
		logger.Error("open dispatchRemoteForward channel", err, map[string]interface{}{
			"module":     "conn",
			"accessId":   fwd.accessId,
			"remoteAddr": remoteAddr,
		})
		facadeConn.Close()
		return
	}
	chId := fwd.chanCounter.Add(1)

	defer func() {
		fwd.chanCounter.Add(-1)
		if value, loaded := fwd.chanMap.LoadAndDelete(chId); loaded {
			value.(*fwdConn).Close()
		}
	}()

	fwd.chanMap.Store(chId, &fwdConn{
		ch:   gosshChan,
		conn: facadeConn,
	})

	go func() {
		defer func() {
			facadeConn.Close()
			gosshChan.Close()
		}()
		_, e := io.Copy(facadeConn, gosshChan)
		if nil != e {
			logger.ErrorN("io.Copy facade write", e)
		}
	}()
	_, e := io.Copy(gosshChan, facadeConn)
	if nil != e {
		logger.ErrorN("io.Copy conn write", e)
	}
}

func (fwd *forwarder) Close() error {
	fwd.cancelFunc()
	fwd.chanMap.Range(func(key, value any) bool {
		value.(*fwdConn).Close()
		return true
	})

	if dbg, ok := fwd.sess.Context().Value(sshDebugServer).(*debugServer); ok {
		dbg.Close()
	}
	logger.DebugN("close all pairs fwd conn")
	return fwd.sess.Exit(0)
}
