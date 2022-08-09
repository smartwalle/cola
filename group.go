package cola

type group struct {
	weight  int32 // 权重
	reject  int   // 拒绝数量
	accept  int   // 接受数量
	actions []*action
}

func newGroup(weight int32) *group {
	var g = &group{}
	g.weight = weight
	g.actions = make([]*action, 0, 4)
	return g
}

func (this *group) push(action *action) {
	this.actions = append(this.actions, action)
}

type GroupList []*group

func (this GroupList) Len() int {
	return len(this)
}

func (this GroupList) Less(i, j int) bool {
	if this[i].weight > this[j].weight {
		return true
	}
	return false
}

func (this GroupList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
