package stat

import (
	"github.com/echogy-io/echogy/pkg/queue"
	"github.com/gliderlabs/ssh"
	"net/http"
)

const (
	requests    = "proxiedRequests"
	requestStat = "requestStat"
)

type Stat struct {
	Receive   int64
	Send      int64
	Request   int
	Response  int
	ConnCount int
	TotalConn int
}

type RequestEntity struct {
	*http.Response
	*http.Request
	UseTime int64
}

func GetQueue(ctx ssh.Context) *queue.FixedQueue {
	q := ctx.Value(requests)
	if nil != q {
		return q.(*queue.FixedQueue)
	} else {
		fq := queue.NewFixedQueue(52)
		ctx.SetValue(requests, fq)
		return fq
	}
}

func GetStat(ctx ssh.Context) *Stat {
	q := ctx.Value(requestStat)
	if nil != q {
		return q.(*Stat)
	} else {
		s := &Stat{}
		ctx.SetValue(requestStat, s)
		return s
	}
}

func Put(ctx ssh.Context, w *http.Response, r *http.Request, t int64) {
	q := GetQueue(ctx)
	s := GetStat(ctx)

	s.Send += w.ContentLength
	s.Receive += r.ContentLength

	s.Request += 1
	s.Response += 1

	q.Push(&RequestEntity{
		Request:  r,
		Response: w,
		UseTime:  t,
	})
}
