package demo

import (
	"fmt"
	"github.com/wxsatellite/goweb/framework"
)

type Service struct {
	IService

	container framework.Container
}

func (*Service) SayHello() {
	fmt.Println("hello")
}

func New(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)

	fmt.Println("new demo service")
	// 返回实例
	return &Service{container: c}, nil
}
