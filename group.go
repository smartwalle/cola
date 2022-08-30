package cola

type group[T any] struct {
	weight  int32 // 权重
	reject  int   // 拒绝数量
	accept  int   // 接受数量
	actions []*action[T]
}

func newGroup[T any](weight int32) *group[T] {
	var g = &group[T]{}
	g.weight = weight
	g.actions = make([]*action[T], 0, 4)
	return g
}

func (this *group[T]) push(action *action[T]) {
	this.actions = append(this.actions, action)
}

type GroupList[T any] []*group[T]

func (this GroupList[T]) Len() int {
	return len(this)
}

func (this GroupList[T]) Less(i, j int) bool {
	if this[i].weight > this[j].weight {
		return true
	}
	return false
}

func (this GroupList[T]) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
