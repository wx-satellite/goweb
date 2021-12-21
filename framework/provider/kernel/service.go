package kernel

import (
	"errors"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/gin"
)

type Service struct {
	engine    *gin.Engine
	container framework.Container
}

func New(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	engine := params[1].(*gin.Engine)
	return &Service{engine: engine, container: container}, nil
}

func (s *Service) Engine() *gin.Engine {
	return s.engine
}
