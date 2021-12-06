package framework

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Context 基于 request 的 Context 进行了封装（ 自定义 Context 结构封装了 net/http 标准库主逻辑流程产生的 Context ）
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter

	// 是否超时标记位
	hasTimeout bool

	// 写保护机制
	writerMux *sync.RWMutex

	// 在中间件注册的回调函数中，只有 framework.Context 这个数据结构作为参数，所以在 Context 中也需要保存这个控制器链路 (handlers)，
	// 并且要记录下当前执行到了哪个控制器（index）
	// 当前请求的 handler 链条
	handlers []ControllerHandler
	// 当前请求调用到链条的哪个节点
	index int

	// 存放解析之后的路由参数
	params map[string]string
}

// NewContext 构造函数
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		hasTimeout:     false,
		writerMux:      &sync.RWMutex{},
		index:          -1,
	}
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}
func (ctx *Context) SetRequest(request *http.Request) {
	ctx.request = request
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}
func (ctx *Context) Params() map[string]string {
	return ctx.params
}

// 执行handler
func (ctx *Context) Next() error {
	ctx.index++
	for ctx.index <= len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// Deadline、Done、Err、Value 是为了实现接口 Context。
// 底层实现逻辑是委派给了 request 的 Context
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// 获取响应实例
func (ctx *Context) Response() http.ResponseWriter {
	return ctx.responseWriter
}

// SetHasTimeout 设置超时
func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
	return
}

// HasTimeout 获取超时的标记
func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// WriterMux 获取读写锁
func (ctx *Context) WriterMux() *sync.RWMutex {
	return ctx.writerMux
}

// QueryKey 获取指定key的GET参数，返回切片，例如：name=ha&name=xi --> []string{"ha","xi"}
func (ctx *Context) QueryKey(key string, def []string) []string {
	params := ctx.QueryAll()
	if val, ok := params[key]; ok {
		return val
	}
	return def
}

// FormKey 获取指定key的POST参数，多个值的时候返回切片
func (ctx *Context) FormKey(key string, def []string) []string {
	res := ctx.FormAll()
	if val, ok := res[key]; !ok {
		return val
	}
	return def
}

// HTML 返回网页信息
func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}
