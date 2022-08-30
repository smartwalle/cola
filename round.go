package cola

import (
	"context"
	"sort"
	"sync"
)

type round[T any] struct {
	mu     *sync.Mutex
	check  chan struct{}
	groups GroupList[T]
	done   bool
}

func newRound[T any]() *round[T] {
	var r = &round[T]{}
	r.mu = &sync.Mutex{}
	r.groups = make(GroupList[T], 0, 12)
	r.done = false
	return r
}

func (this *round[T]) finish(weight int32, status int32) {
	this.mu.Lock()
	defer this.mu.Unlock()

	var nGroup *group[T]
	for _, g := range this.groups {
		if g.weight == weight {
			nGroup = g
			break
		}
	}

	if nGroup != nil {
		if status == statusAccept {
			nGroup.accept++
		} else if status == statusReject {
			nGroup.reject++
		}
	}

	select {
	case this.check <- struct{}{}:
	default:
	}
}

func (this *round[T]) add(action *action[T]) {
	if action == nil {
		return
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	var nGroup *group[T]
	for _, g := range this.groups {
		if g.weight == action.weight {
			nGroup = g
			break
		}
	}

	if nGroup == nil {
		nGroup = newGroup[T](action.weight)
		this.groups = append(this.groups, nGroup)

		sort.Sort(this.groups)
	}

	nGroup.push(action)
}

func (this *round[T]) tick(ctx context.Context, finished func([]T)) {
	this.mu.Lock()
	var total = cap(this.groups)
	this.check = make(chan struct{}, total)
	this.mu.Unlock()

	var victors []T

	defer func() {
		this.done = true
		close(this.check)
		if finished != nil {
			finished(victors)
		}
	}()

	if ok, result := this.exec(false); ok {
		victors = result
		return
	}

	for {
		select {
		case <-ctx.Done():
			_, victors = this.exec(true)
			return
		case <-this.check:
			if ok, result := this.exec(false); ok {
				victors = result
				return
			}
		}
	}
}

func (this *round[T]) exec(focus bool) (bool, []T) {
	this.mu.Lock()
	defer this.mu.Unlock()

	var done = false

	for _, g := range this.groups {
		var finish = g.accept + g.reject
		var total = len(g.actions)

		// 所有决策是否已经完成
		done = finish == total

		// 如果决策未完成并且不是强制要求出结果，则直接返回
		if done == false && focus == false {
			return false, nil
		}

		// 1、该组已做出所有决策，并且通过数量大于 0，则表示已决策出结果
		// 2、强制要求出结果，并且通过数量大于 0，则表示已决策出结果
		if g.accept > 0 && (done || focus) {
			var victors = make([]T, 0, g.accept)
			for _, m := range g.actions {
				if m.exec() {
					victors = append(victors, m.Data())
				}
			}
			return true, victors
		}
	}
	return done, nil
}
