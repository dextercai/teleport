package main

import (
	tp "github.com/henrylee2cn/teleport"
)

//go:generate go build $GOFILE

func main() {
	defer tp.FlushLogger()
	srv := tp.NewPeer(
		tp.PeerConfig{
			PrintDetail: false,
			ListenPort:  9090,
		},
		new(earlyResult),
	)
	srv.ListenAndServe()
}

type earlyResult struct{}

func (e *earlyResult) Name() string {
	return "early_result"
}

func (e *earlyResult) PostAccept(sess tp.PreSession) *tp.Status {
	var rigthServiceMethod bool
	input, stat := sess.Receive(func(header tp.Header) interface{} {
		if header.ServiceMethod() == "/early/ping" {
			rigthServiceMethod = true
			return new(map[string]string)
		}
		return nil
	})
	if !stat.OK() {
		return stat
	}

	var result string
	if !rigthServiceMethod {
		stat = tp.NewStatus(10005, "unexpected request", "")
	} else {
		body := *input.Body().(*map[string]string)
		if body["author"] != "henrylee2cn" {
			stat = tp.NewStatus(10005, "incorrect author", body["author"])
		} else {
			stat = nil
			result = "OK"
		}
	}
	return sess.Send(
		"/early/pong",
		result,
		stat,
	)
}
