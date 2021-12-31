package log

import (
	"github.com/pkg/errors"
	"github.com/wxsatellite/goweb/framework"
	"os"
)

// ConsoleService 控制台输出日志
type ConsoleService struct {

	// BaseService 实现了 Log 接口，匿名嵌套之后，ConsoleService 自然而然也实现了 Log 接口
	// ConsoleService以及其他Service 相较于 BaseService 就是 outPut 不一样，因此其他逻辑都抽出来放在 BaseService 中复用，
	// 然后提供 outPut 的设置函数 SetOutput 来满足不同的 Service。（ 这个实现思路可以学习一下！ ）
	BaseService
}

func NewConsoleService(params ...interface{}) (interface{}, error) {
	if len(params) < 4 {
		return nil, errors.New("params error")
	}
	container := params[0].(framework.Container)
	level := params[1].(Level)
	ctxFielder := params[2].(CtxFielder)
	formatter := params[3].(Formatter)

	service := &ConsoleService{}
	service.SetCtxFielder(ctxFielder)
	service.SetFormatter(formatter)
	service.SetLevel(level)
	service.SetOutput(os.Stdout)
	service.container = container
	return service, nil
}
