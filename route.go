package main

import "goweb/framework"

func registerRouter(core *framework.Core) {
	core.Get("/foo", FooControllerHandler)
}
