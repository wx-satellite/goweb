package env

import (
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/app"
)

type Provider struct {
	Folder string
}

func (p *Provider) Boot(container framework.Container) (err error) {
	// 存在就不设置了
	if p.Folder != "" {
		return
	}
	appService := container.MustMake(app.Key).(app.App)
	p.Folder = appService.BaseFolder()
	return
}

func (p *Provider) Params(container framework.Container) []interface{} {
	return []interface{}{container, p.Folder}
}

func (p *Provider) Name() string {
	return Key
}

func (p *Provider) IsDefer() bool {
	return false
}

func (p *Provider) Register(container framework.Container) framework.NewInstance {
	return New
}
