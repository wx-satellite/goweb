package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
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
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
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

// GetRequest 获取请求实例
func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

// 获取响应实例
func (ctx *Context) GetResponse() http.ResponseWriter {
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

// QueryInt 获取指定key的GET参数的int类型，以最后一个为准，例如：age=12&age=24 --> 24
func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def
	}
	intVal, err := strconv.Atoi(val[len(val)-1])
	if err != nil {
		return def
	}
	return intVal
}

// QueryString 获取指定key的GET参数，以最后一个为准，例如：name=ha&name=xi --> xi
func (ctx *Context) QueryString(key, def string) string {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def
	}
	return val[len(val)-1]
}

// QueryKey 获取指定key的GET参数，返回切片，例如：name=ha&name=xi --> []string{"ha","xi"}
func (ctx *Context) QueryKey(key string, def []string) []string {
	params := ctx.QueryAll()
	if val, ok := params[key]; ok {
		return val
	}
	return def
}

// QueryAll 获取所有GET参数
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}
	return ctx.request.URL.Query()
}

// FormInt 获取指定key的POST参数，多个值的时候只返回最后一个，并转成数字类型
func (ctx *Context) FormInt(key string, def int) int {
	res := ctx.FormAll()
	val, ok := res[key]
	if !ok || len(val) <= 0 {
		return def
	}
	intVal, err := strconv.Atoi(val[len(val)-1])
	if err != nil {
		return def
	}
	return intVal
}

// FormString 获取指定key的POST参数，多个值的时候只返回最后一个
func (ctx *Context) FormString(key, def string) string {
	res := ctx.FormAll()
	val, ok := res[key]
	if !ok || len(val) <= 0 {
		return def
	}
	return val[len(val)-1]

}

// FormKey 获取指定key的POST参数，多个值的时候返回切片
func (ctx *Context) FormKey(key string, def []string) []string {
	res := ctx.FormAll()
	if val, ok := res[key]; !ok {
		return val
	}
	return def
}

// FormAll 获取所有的POST参数
func (ctx *Context) FormAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}
	// Form：存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。
	// PostForm：存储了 post、put 参数，在使用之前需要调用 ParseForm 方法。
	_ = ctx.request.ParseForm()
	return ctx.request.PostForm
}

// BindJson 绑定json数据到结构体中
func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request == nil {
		return errors.New("ctx.request body empty")
	}
	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return err
	}
	// io 只能读取一次，所以这里为了不影响后续Body的读取使用NopCloser
	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = json.Unmarshal(body, obj)
	return err
}

// Json 返回json格式的响应体
func (ctx *Context) Json(status int, obj interface{}) error {
	// 当超时的时候会往response中写入数据，并且设置 hasTimeout true
	// 如果 hasTimeout 为 true 则表示已经向 response 写入了数据，这里就直接return
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = ctx.responseWriter.Write(byt)
	return err
}

// HTML 返回网页信息
func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

// Text 返回文本信息
func (ctx *Context) Text(status int, obj string) error {
	return nil
}
