package main

import "github.com/wxsatellite/goweb/framework"

func registerRouter(core *framework.Core) {
	core.Get("/foo", FooControllerHandler)
}
