package framework

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Core 框架核心结构
type Core struct {
	// 二级map，主要处理静态路由，因为支持rest风格，所以还要根据请求方法分类（ http方法和静态路由匹配 ）
	//router map[string]map[string]ControllerHandler
	// 实际在实现的过程中，其实第一层是request-uri第二层是method也是可以的，或者直接一层，将request-uri+method做为key

	// router 为了支持动态路由，将二级map替换成字典树
	router map[string]*Tree
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	//getRouter := map[string]ControllerHandler{}
	//postRouter := map[string]ControllerHandler{}
	//putRouter := map[string]ControllerHandler{}
	//deleteRouter := map[string]ControllerHandler{}
	//router := map[string]map[string]ControllerHandler{}
	//router["GET"] = getRouter
	//router["POST"] = postRouter
	//router["PUT"] = putRouter
	//router["DELETE"] = deleteRouter

	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{
		router: router,
	}
}

// Group 路由分组，批量通用前缀
func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// Get get请求的路由注册
func (c *Core) Get(uri string, handler ControllerHandler) {
	// 注册的时候将URL全部大写，在匹配的时候也需要转成大写匹配。这样子实现的路由就是"大小写不敏感"的，对使用者的容错率增加
	//upperUri := strings.ToUpper(uri)
	//c.router["GET"][upperUri] = handler
	fmt.Println(uri)
	if err := c.router["GET"].AddRouter(uri, handler); err != nil {
		log.Fatal("add router error：", err)
	}
}

// Post post请求的路由注册
func (c *Core) Post(uri string, handler ControllerHandler) {
	if err := c.router["POST"].AddRouter(uri, handler); err != nil {
		log.Fatal("add router error：", err)
	}
}

// Put put请求的路由注册
func (c *Core) Put(uri string, handler ControllerHandler) {
	if err := c.router["PUT"].AddRouter(uri, handler); err != nil {
		log.Fatal("add router error：", err)
	}
}

// Delete delete请求的路由注册
func (c *Core) Delete(uri string, handler ControllerHandler) {
	if err := c.router["DELETE"].AddRouter(uri, handler); err != nil {
		log.Fatal("add router error：", err)
	}
}

// MatchRouter 匹配路由
func (c *Core) MatchRouter(request *http.Request) ControllerHandler {
	upperMethod := strings.ToUpper(request.Method)
	uri := request.URL.Path

	tree, ok := c.router[upperMethod]
	if !ok {
		return nil
	}
	return tree.FindHandler(uri)
}

// ServeHTTP 实现 Handler 接口
func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("core.ServeHTTP")

	ctx := NewContext(w, r)

	handler := c.MatchRouter(r)
	if handler == nil {
		_ = ctx.Json(404, "not found")
		return
	}
	err := handler(ctx)
	if err != nil {
		_ = ctx.Json(500, "inner error")
		return
	}
}
