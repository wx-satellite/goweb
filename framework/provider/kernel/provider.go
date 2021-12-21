package kernel

import (
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/gin"
)

type Provider struct {
	// 这个服务提供者可以在注册服务的时候传递 Web 引擎，如果没有传递，则需要在启动的时候默认初始化。
	Engine *gin.Engine
}

func (*Provider) Register(container framework.Container) framework.NewInstance {
	return New
}

func (p *Provider) Boot(container framework.Container) error {
	if p.Engine == nil {
		p.Engine = gin.Default()
	}
	// engine 创建的时候其实会初始化container容器，这里需要进行覆盖
	p.Engine.SetContainer(container)
	return nil
}

func (p *Provider) Params(container framework.Container) []interface{} {
	return []interface{}{container, p.Engine}
}

func (*Provider) IsDefer() bool {
	return false
}

func (*Provider) Name() string {
	return Key
}
