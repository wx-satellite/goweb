package main

import (
	"context"
	"goweb/framework"
	"log"
	"time"
)

func FooControllerHandler(ctx *framework.Context) error {

	// 生成超时的context
	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), 1*time.Second)
	defer cancel()

	// 创建一个新的 goroutine 来处理业务逻辑
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// goroutine 的异常都是需要自己捕获的，不存在父goroutine捕获子goroutine
				// panicChan 就是告知父goroutine，"我"异常了
				panicChan <- p
			}
		}()
		// time.Sleep 来模拟具体业务逻辑的处理时间
		time.Sleep(10 * time.Second)
		_ = ctx.Json(200)
		// 通过 finish 通道告知父goroutine处理结束
		finish <- struct{}{}
	}()

	select {
	// 异常事件、超时事件触发时，需要往 responseWriter 中写入信息，这个时候不再允许其他 Goroutine 操作 responseWriter
	// 否则会导致 responseWriter 中的信息出现乱序。解决方案：锁
	case p := <-panicChan: // 监听panic
		log.Println("请求异常了", p)
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		_ = ctx.Json(500)
	case <-finish: // 监听结束
		_ = ctx.Json(200)
	case <-durationCtx.Done(): // 监听超时
		log.Println("请求超时了")
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		_ = ctx.Json(500)
		ctx.SetHasTimeout()
	}
	return nil
}
