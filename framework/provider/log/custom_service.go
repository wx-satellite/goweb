package log

import (
	"github.com/pkg/errors"
	"github.com/wxsatellite/goweb/framework"
	"io"
)

// CustomService 自定义日志输出
type CustomService struct {
	BaseService
}

func NewCustomService(params ...interface{}) (interface{}, error) {
	if len(params) < 5 {
		return nil, errors.New("params error")
	}
	c := params[0].(framework.Container)
	level := params[1].(Level)
	ctxFielder := params[2].(CtxFielder)
	formatter := params[3].(Formatter)
	output := params[4].(io.Writer)

	log := &CustomService{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(output)
	log.container = c
	return log, nil
}
