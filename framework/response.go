package framework

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
