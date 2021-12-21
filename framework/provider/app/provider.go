package app

import "github.com/wxsatellite/goweb/framework"

type Provider struct {
	BaseFolder string
}

func (*Provider) Register(container framework.Container) framework.NewInstance {
	return New
}

func (p *Provider) Params(container framework.Container) []interface{} {
	return []interface{}{container, p.BaseFolder}
}

func (*Provider) IsDefer() bool {
	return false
}

func (*Provider) Boot(container framework.Container) error {
	return nil
}

func (*Provider) Name() string {
	return Key
}
