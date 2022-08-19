package cola

import (
	"context"
	"github.com/smartwalle/task"
	"sync"
	"time"
)

type Manager[T any] struct {
	mu    *sync.Mutex
	round *round[T]
	task  task.Manager
}

func New[T any]() *Manager[T] {
	var m = &Manager[T]{}
	m.mu = &sync.Mutex{}
	m.task = task.New(task.WithWorker(2))
	m.task.Run()
	return m
}

func (this *Manager[T]) Add(data T, weight int32, handler func(data T)) Action[T] {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.round == nil {
		this.round = newRound[T]()
	}
	var nAction = newAction[T](data, weight, this.round, handler)
	this.round.add(nAction)
	return nAction
}

func (this *Manager[T]) Tick(timeout time.Duration, finished func(victors []T), opts ...TickOption) {
	var ctx, _ = context.WithTimeout(context.Background(), timeout)
	this.tick(ctx, finished, opts...)
}

func (this *Manager[T]) TickWithDeadline(deadline time.Time, finished func(victors []T), opts ...TickOption) {
	var ctx, _ = context.WithDeadline(context.Background(), deadline)
	this.tick(ctx, finished, opts...)
}

func (this *Manager[T]) Close() {
	this.task.Close()
}

func (this *Manager[T]) tick(ctx context.Context, finished func([]T), opts ...TickOption) {
	this.mu.Lock()
	var current = this.round
	this.round = nil
	this.mu.Unlock()

	if current != nil {
		var nOpts = &TickOptions{}

		for _, opt := range opts {
			if opt != nil {
				opt(nOpts)
			}
		}

		if nOpts.waiter != nil {
			nOpts.waiter.Add(1)
		}

		this.task.AddTask(func(arg interface{}) {
			current.tick(ctx, finished, nOpts)
		})
	}
}

type TickOptions struct {
	waiter Waiter
}

type TickOption func(opts *TickOptions)

func WithWaiter(waiter Waiter) TickOption {
	return func(opts *TickOptions) {
		opts.waiter = waiter
	}
}
