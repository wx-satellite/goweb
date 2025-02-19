package middleware

import (
	"context"
	"fmt"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/gin"
	"log"
	"time"
)

// TimeoutHandler Pipeline 的方式改造中间件的实现：一层层嵌套不好用，如果我们将每个核心控制器所需要的中间件，使用一个数组链接（Chain）起来，形成一条流水线（Pipeline）
// 因为实际业务逻辑的控制器和中间件的类型相同即为 ControllerHandler，因此可以将它们成一个 ControllerHandler数组，也就是控制器链
func TimeoutHandler(d time.Duration) framework.ControllerHandler {
	return func(ctx *gin.Context) error {
		durationCtx, cancel := context.WithTimeout(ctx.Request.Context(), d)
		defer cancel()

		// 将新的context设置到request中
		ctx.Request = ctx.Request.Clone(durationCtx)

		finishChan := make(chan struct{}, 1)
		panicChan := make(chan struct{}, 1)

		go func() {
			if p := recover(); p != nil {
				panicChan <- struct{}{}
			}
			ctx.Next()

			finishChan <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			log.Printf("request error：%v", p)
			_ = ctx.IJson(500)
		case <-finishChan:
		case <-durationCtx.Done():
			log.Printf("request timeout")
			_ = ctx.IJson(500)
		}

		return nil
	}
}

// TimeoutHandlerV1 装饰器的实现方式。类似于洋葱。
// 这种实现方式有两个不好：
//  中间件是循环嵌套的，如果有多个中间件的时候，嵌套会很长，不优雅： TimeoutHandler(LogHandler(recoveryHandler(UserLoginController)))
//  只能为单个业务控制器设置中间件，不能批量设置
func TimeoutHandlerV1(f framework.ControllerHandler, d time.Duration) framework.ControllerHandler {
	return func(ctx *gin.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan struct{}, 1)

		// 设置超时
		durationCtx, cancel := context.WithTimeout(ctx.Request.Context(), d)
		defer cancel()

		// 将新的context设置到request中
		ctx.Request = ctx.Request.Clone(durationCtx)

		// 处理业务逻辑
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- struct{}{}
				}
			}()
			// 实际的业务逻辑
			err := f(ctx)
			log.Println(err)

			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			log.Println(p)
			_ = ctx.IJson(500)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			_ = ctx.IJson(500)
		}
		return nil
	}
}
