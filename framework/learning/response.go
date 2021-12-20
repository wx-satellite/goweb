package learning

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// IResponse 很多方法返回的是接口本身，这么设计允许使用者进行链式调用。
// 链式调用的好处能很大提升代码的阅读性： c.SetOkStatus().Json("ok, UserLoginController: " + foo)
type IResponse interface {

	// Json json的形式输出
	Json(obj interface{}) IResponse

	// Jsonp jsonp的形式输出
	Jsonp(obj interface{}) IResponse

	// Xml xml的形式输出
	Xml(obj interface{}) IResponse

	// Html html输出
	Html(template string, obj interface{}) IResponse

	// Text 文本的形式输出
	Text(format string, values ...interface{}) IResponse

	// Redirect 重定向
	Redirect(path string) IResponse

	// SetHeader 设置响应头
	SetHeader(key string, val string) IResponse

	// SetCookie 设置cookie
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	// SetStatus 设置状态码
	SetStatus(code int) IResponse

	// SetOkStatus 设置成功的状态码
	SetOkStatus() IResponse
}

func (ctx *Context) SetStatus(code int) IResponse {
	ctx.Response().WriteHeader(code)
	return ctx
}

func (ctx *Context) SetOkStatus() IResponse {
	ctx.Response().WriteHeader(http.StatusOK)
	return ctx
}

func (ctx *Context) SetHeader(key, val string) IResponse {
	ctx.Response().Header().Set(key, val)
	return ctx
}

func (ctx *Context) Redirect(path string) IResponse {
	http.Redirect(ctx.Response(), ctx.Request(), path, http.StatusFound)
	return ctx
}

func (ctx *Context) Xml(obj interface{}) IResponse {
	bs, err := xml.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/xml; charset=utf-8")
	_, _ = ctx.Response().Write(bs)
	return ctx
}

func (ctx *Context) Json(obj interface{}) IResponse {
	bs, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/json;charset=utf-8")
	_, _ = ctx.Response().Write(bs)
	return ctx
}
func (ctx *Context) Html(file string, obj interface{}) IResponse {
	paths := strings.Split(file, "/")
	if len(paths) <= 0 {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	lastPath := paths[len(paths)-1]
	name := lastPath[:strings.Index(lastPath, ".")]

	// 渲染html：模版+数据

	// 创建模版
	tmpl, _ := template.New(name).ParseFiles(file)
	// 填充数据
	if err := tmpl.Execute(ctx.Response(), obj); err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/html;charset=utf-8")
	return ctx
}

func (ctx *Context) Text(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.SetHeader("Content-Type", "application/text;charset=utf-8")
	_, _ = ctx.Response().Write([]byte(out))
	return ctx
}

func (ctx *Context) SetCookie(key string, val string, maxAge int, path string, domain string, secure bool, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.Response(), &http.Cookie{
		Name:     key,
		Value:    val,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: 1,
	})
	return ctx
}

// Jsonp 浏览器的同源策略导致A域名的页面不能通过ajax获取B域名的数据，解决方法是通过 script 标签（不受同源策略的限制），返回一段js代码
func (ctx *Context) Jsonp(obj interface{}) IResponse {
	callbackFunc, _ := ctx.QueryString("callback", "callback_function")
	ctx.SetHeader("Content-Type", "application/javascript")

	//输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成 XSS 攻击
	// callbackFunc 是前端传递的，有可能包含js代码
	callbackFunc = template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := ctx.Response().Write([]byte(callbackFunc))
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}

	// 输出左括号
	_, err = ctx.Response().Write([]byte("("))
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}

	// 输出数据
	bs, _ := json.Marshal(obj)
	_, err = ctx.Response().Write(bs)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}

	// 输出右括号
	_, err = ctx.Response().Write([]byte(")"))
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	return ctx
}
