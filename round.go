package cola

import (
	"sort"
	"sync"
)

type round struct {
	mu     *sync.Mutex
	groups GroupList
	check  chan struct{}
	done   bool
}

func newRound() *round {
	var r = &round{}
	r.mu = &sync.Mutex{}
	r.groups = make(GroupList, 0, 12)
	r.done = false
	return r
}

func (this *round) finish(weight int32, status int32) {
	this.mu.Lock()
	defer this.mu.Unlock()

	var nGroup *group
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

func (this *round) add(action *action) {
	if action == nil {
		return
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	var nGroup *group
	for _, g := range this.groups {
		if g.weight == action.weight {
			nGroup = g
			break
		}
	}

	if nGroup == nil {
		nGroup = newGroup(action.weight)
		this.groups = append(this.groups, nGroup)

		sort.Sort(this.groups)
	}

	nGroup.push(action)
}

func (this *round) tick(opt *tickOption) {
	this.mu.Lock()
	var total = cap(this.groups)
	this.check = make(chan struct{}, total)
	this.mu.Unlock()

	defer func() {
		this.done = true
		close(this.check)
		if opt.finish != nil {
			opt.finish()
		}
		if opt.waiter != nil {
			opt.waiter.Done()
		}
	}()

	if this.exec(false) {
		return
	}

	for {
		select {
		case <-opt.context.Done():
			this.exec(true)
			return
		case <-this.check:
			if this.exec(false) {
				return
			}
		}
	}
}

func (this *round) exec(focus bool) bool {
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
			return false
		}

		// 1、该组已做出所有决策，并且通过数量大于 0，则表示已决策出结果
		// 2、强制要求出结果，并且通过数量大于 0，则表示已决策出结果
		if g.accept > 0 && (done || focus) {
			for _, m := range g.actions {
				m.exec()
			}
			return true
		}
	}
	return done
}
