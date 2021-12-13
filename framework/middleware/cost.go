package middleware

import (
	"github.com/wxsatellite/goweb/framework"
	"log"
	"time"
)

// Cost 统计请求消耗的时间
func Cost() framework.ControllerHandler {
	return func(ctx *framework.Context) error {

		// 记录开始时间
		startTime := time.Now()

		// 执行业务逻辑
		_ = ctx.Next()

		// 记录结束时间
		endTime := time.Now()

		// 消耗的时间
		cost := endTime.Sub(startTime)

		// 记录日志
		log.Printf("api uri：%v，cost：%v\n", ctx.Request().RequestURI, cost.Seconds())

		return nil
	}
}
