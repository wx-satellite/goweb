package framework

import "github.com/wxsatellite/goweb/framework/gin"

/**
type ControllerHandler func(ctx *gin.Context) error

gin 的 handler 定义：
type HandlerFunc func(*Context)

相比之后，项目中的handler多了一个 error 错误。gin 的作者认为中断一个请求返回一个 error 并没有什么用。
他希望中断一个请求时操作 response，比如设置状态码、设置返回错误信息等等，而不是希望用 error 来进行返回。

*/
type ControllerHandler func(ctx *gin.Context) error
