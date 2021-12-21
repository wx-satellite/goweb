package gin

import (
	"context"
	"github.com/wxsatellite/goweb/framework"
)

// Engine 负责在容器中绑定服务提供者，Context 负责从容器中获取服务提供者

func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

func (engine *Engine) SetContainer(container framework.Container) {
	engine.container = container
}

func (c *Context) Make(key string) (interface{}, error) {
	return c.container.Make(key)
}

func (c *Context) MustMake(key string) interface{} {
	return c.container.MustMake(key)
}

func (c *Context) MakeNew(key string, params []interface{}) (interface{}, error) {
	return c.container.MakeNew(key, params)
}

func (c *Context) BaseContext() context.Context {
	return c.Request.Context()
}
