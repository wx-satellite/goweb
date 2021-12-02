package middleware

import (
	"goweb/framework"
	"log"
)

// Recovery 捕获异常
func Recovery() framework.ControllerHandler {
	return func(ctx *framework.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic：%v", err)
				_ = ctx.Json(500, "server error")
			}
		}()
		_ = ctx.Next()
		return nil
	}
}
