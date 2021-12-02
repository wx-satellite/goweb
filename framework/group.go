package framework

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)

	// 实现嵌套 group
	// eg：core.Group("/user").Group("/name")
	Group(string) IGroup
}

// Group 这里不是直接使用 Group 结构而是使用 IGroup 接口做抽象是因为如果后续
// 我们发现设计的 Group 结构不满足需求了，需要引入另一个 Group2 结构，为了保证
// 调用方受影响程度最少，所以做了接口抽象
type Group struct {
	core   *Core
	prefix string
}

// NewGroup 实例化 group
func NewGroup(core *Core, prefix string) IGroup {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

// Group 支持 group 嵌套
func (g *Group) Group(prefix string) IGroup {
	return &Group{core: g.core, prefix: g.prefix + prefix}
}

// Get 设置get路由
func (g *Group) Get(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Get(uri, handler)
	return
}

// Post 设置post路由
func (g *Group) Post(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Post(uri, handler)
	return
}

// Put 设置put路由
func (g *Group) Put(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Put(uri, handler)
	return
}

// Delete 设置delete路由
func (g *Group) Delete(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.core.Delete(uri, handler)
	return
}
