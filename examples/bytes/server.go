package main

import (
	"time"

	tp "github.com/henrylee2cn/teleport"
)

//go:generate go build $GOFILE

func main() {
	defer tp.FlushLogger()
	go tp.GraceSignal()
	tp.SetShutdown(time.Second*20, nil, nil)
	var peer = tp.NewPeer(tp.PeerConfig{
		SlowCometDuration: time.Millisecond * 500,
		PrintDetail:       true,
		ListenPort:        9090,
	})
	group := peer.SubRoute("group")
	group.RouteCall(new(Home))
	peer.SetUnknownCall(UnknownCallHandle)
	peer.ListenAndServe()
}

// Home controller
type Home struct {
	tp.CallCtx
}

// Test handler
func (h *Home) Test(arg *[]byte) ([]byte, *tp.Status) {
	h.Session().Push("/push/test", []byte("test push text"))
	tp.Debugf("HomeCallHandle: codec: %d, arg: %s", h.GetBodyCodec(), *arg)
	return []byte("test call result text"), nil
}

// UnknownCallHandle handles unknown call message
func UnknownCallHandle(ctx tp.UnknownCallCtx) (interface{}, *tp.Status) {
	ctx.Session().Push("/push/test", []byte("test unknown push text"))
	var arg []byte
	codecID, err := ctx.Bind(&arg)
	if err != nil {
		return nil, tp.NewStatus(1001, "bind error", err.Error())
	}
	tp.Debugf("UnknownCallHandle: codec: %d, arg: %s", codecID, arg)
	return []byte("test unknown call result text"), nil
}
