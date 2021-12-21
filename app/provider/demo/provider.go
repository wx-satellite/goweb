package demo

import (
	"fmt"
	"github.com/wxsatellite/goweb/framework"
)

type ServiceProvider struct {
}

func (*ServiceProvider) Name() string {
	return Key
}

func (*ServiceProvider) Register(container framework.Container) framework.NewInstance {
	return New
}

func (*ServiceProvider) IsDefer() bool {
	return true
}

func (*ServiceProvider) Boot(container framework.Container) error {
	fmt.Println("demo service boot")

	return nil
}

func (*ServiceProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}
