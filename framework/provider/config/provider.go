package config

import (
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/app"
	"github.com/wxsatellite/goweb/framework/provider/env"
)

type Provider struct {
	folder  string
	env     string
	envMaps map[string]string
}

func (p *Provider) Register(container framework.Container) framework.NewInstance {
	return New
}

func (p *Provider) Boot(container framework.Container) (err error) {
	appService := container.MustMake(app.Key).(app.App)
	envService := container.MustMake(env.Key).(env.Env)
	p.folder = appService.ConfigFolder()
	p.env = envService.AppEnv()
	p.envMaps = envService.All()
	return
}

func (p *Provider) Params(container framework.Container) []interface{} {
	return []interface{}{container, p.folder, p.env, p.envMaps}
}

func (p *Provider) Name() string {
	return Key
}

func (p *Provider) IsDefer() bool {
	return false
}
