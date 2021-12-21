package main

import (
	"github.com/wxsatellite/goweb/framework/gin"
)

func registerRouter(core *gin.Engine) {
	core.GET("/foo", FooControllerHandler)

}
