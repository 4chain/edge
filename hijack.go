package echogy

import (
	"bufio"
	"bytes"
	q "github.com/echogy-io/echogy/pkg/queue"
	"net"
	"net/http"
	"time"
)

type Dispatch func(*http.Response, *http.Request, int64)

type hijackHttp struct {
	net.Conn
	hasFwdReq bool
	dispatch  Dispatch
	queue     *q.SyncQueue
}

type request struct {
	*http.Request
	startTime int64
}

func newHijackConn(conn net.Conn) *hijackHttp {
	queue := q.NewSyncQueue(4)
	return &hijackHttp{
		Conn:      conn,
		hasFwdReq: true,
		queue:     queue,
	}
}

func (h *hijackHttp) Read(b []byte) (n int, err error) {
	n, err = bufio.NewReader(h.Conn).Read(b)
	if err != nil {
		return n, err
	}

	reader := bytes.NewReader(b[:n])
	if req, err := http.ReadRequest(bufio.NewReader(reader)); err == nil {
		// add req to queue
		h.queue.Push(&request{
			Request:   req,
			startTime: time.Now().UnixMilli(),
		})
	}
	return n, nil
}

func (h *hijackHttp) SetDispatch(d Dispatch) {
	h.dispatch = d
}

func (h *hijackHttp) Write(b []byte) (n int, err error) {
	n, err = h.Conn.Write(b)
	if nil != err {
		return n, err
	}

	pop := h.queue.Pop()
	if nil != pop {
		r := pop.(*request)
		// pop request from request queue and Try to parse HTTP response from the data
		if nil != h.dispatch {
			reader := bytes.NewReader(b)
			if resp, err := http.ReadResponse(bufio.NewReader(reader), r.Request); nil == err {
				useTime := time.Now().UnixMilli() - r.startTime
				h.dispatch(resp, r.Request, useTime)
			}
		}
	}
	return n, nil
}
