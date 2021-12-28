package distributed

import (
	"github.com/wxsatellite/goweb/framework"
)

type LocalProvider struct {
}

func (p *LocalProvider) Register(container framework.Container) framework.NewInstance {
	return NewLocalProvider
}

func (p *LocalProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (p *LocalProvider) Boot(container framework.Container) error {
	return nil
}

func (p *LocalProvider) IsDefer() bool {
	return false
}

func (p *LocalProvider) Name() string {
	return Key
}
