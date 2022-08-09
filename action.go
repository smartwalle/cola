package cola

import (
	"sync/atomic"
)

type Action interface {
	Key() string

	Reject()

	Accept()

	Valid() bool
}

const (
	statusReject  = 1 // 拒绝
	statusDefault = 2 // 默认，初始化状态
	statusAccept  = 3 // 接受
)

type action struct {
	status  int32
	key     string
	weight  int32
	round   *round
	handler func(key string)
}

func newAction(key string, weight int32, round *round, handler func(key string)) *action {
	var a = &action{}
	a.status = statusDefault
	a.key = key
	a.weight = weight
	a.round = round
	a.handler = handler
	return a
}

func (this *action) Key() string {
	return this.key
}

func (this *action) Reject() {
	if atomic.CompareAndSwapInt32(&this.status, statusDefault, statusReject) {
		this.round.finish(this.weight, statusReject)
	}
}

func (this *action) Accept() {
	if atomic.CompareAndSwapInt32(&this.status, statusDefault, statusAccept) {
		this.round.finish(this.weight, statusAccept)
	}
}

func (this *action) Valid() bool {
	return this.round.done == false
}

func (this *action) exec() {
	if this.handler != nil && atomic.LoadInt32(&this.status) == statusAccept {
		this.handler(this.key)
	}
}
