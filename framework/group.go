package framework

type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	// 实现嵌套 group
	// eg：core.Group("/user").Group("/name")
	Group(string) IGroup

	// 设置中间件
	Use(middlewares ...ControllerHandler)
}

// Group 这里不是直接使用 Group 结构而是使用 IGroup 接口做抽象是因为如果后续
// 我们发现设计的 Group 结构不满足需求了，需要引入另一个 Group2 结构，为了保证
// 调用方受影响程度最少，所以做了接口抽象
type Group struct {
	core        *Core
	prefix      string
	middlewares []ControllerHandler
}

// NewGroup 实例化 group
func NewGroup(core *Core, prefix string) IGroup {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Group 支持 group 嵌套
func (g *Group) Group(prefix string) IGroup {
	return &Group{core: g.core, prefix: g.prefix + prefix, middlewares: g.middlewares}
}

// combineHandlers 合并handlers，使用新创建的切片。如果直接append，例如：handlers := append(g.middlewares, handler...)
// 由于切片底层数组共享的问题可能会导致处理具柄互相覆盖：
/*
s := make([]int64, 0, 4)
s = append(s, 1, 2) // 追加中间件
s1 := append(s, 3) // 追加业务句柄1
fmt.Println(s)
fmt.Println(s1)
s2 := append(s, 4) // 追加业务句柄2
fmt.Println(s)
fmt.Println(s1) // 业务句柄2覆盖了业务句柄1
fmt.Println(s2)
*/
func (g *Group) combineHandlers(handler ...ControllerHandler) []ControllerHandler {
	finalSize := len(g.middlewares) + len(handler)
	mergedHandlers := make([]ControllerHandler, finalSize)
	copy(mergedHandlers, g.middlewares)
	copy(mergedHandlers[len(g.middlewares):], handler)
	return mergedHandlers
}

// Get 设置get路由
func (g *Group) Get(uri string, handler ...ControllerHandler) {
	uri = g.prefix + uri
	g.core.Get(uri, g.combineHandlers(handler...)...)
	return
}

// Post 设置post路由
func (g *Group) Post(uri string, handler ...ControllerHandler) {
	uri = g.prefix + uri
	g.core.Post(uri, g.combineHandlers(handler...)...)
	return
}

// Put 设置put路由
func (g *Group) Put(uri string, handler ...ControllerHandler) {
	uri = g.prefix + uri
	g.core.Put(uri, g.combineHandlers(handler...)...)
	return
}

// Delete 设置delete路由
func (g *Group) Delete(uri string, handler ...ControllerHandler) {
	uri = g.prefix + uri
	g.core.Delete(uri, g.combineHandlers(handler...)...)
	return
}
