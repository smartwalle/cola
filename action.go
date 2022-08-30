package cola

import (
	"sync/atomic"
)

type Action[T any] interface {
	Data() T

	Reject()

	Accept()

	Valid() bool
}

const (
	statusReject  = 1 // 拒绝
	statusDefault = 2 // 默认，初始化状态
	statusAccept  = 3 // 接受
)

type action[T any] struct {
	data    T
	round   *round[T]
	handler func(T)
	status  int32
	weight  int32
}

func newAction[T any](data T, weight int32, round *round[T], handler func(data T)) *action[T] {
	var a = &action[T]{}
	a.status = statusDefault
	a.data = data
	a.weight = weight
	a.round = round
	a.handler = handler
	return a
}

func (this *action[T]) Data() T {
	return this.data
}

func (this *action[T]) Reject() {
	if atomic.CompareAndSwapInt32(&this.status, statusDefault, statusReject) {
		this.round.finish(this.weight, statusReject)
	}
}

func (this *action[T]) Accept() {
	if atomic.CompareAndSwapInt32(&this.status, statusDefault, statusAccept) {
		this.round.finish(this.weight, statusAccept)
	}
}

func (this *action[T]) Valid() bool {
	return this.round.done == false
}

func (this *action[T]) exec() bool {
	if atomic.LoadInt32(&this.status) == statusAccept {
		if this.handler != nil {
			this.handler(this.data)
		}
		return true
	}
	return false
}
