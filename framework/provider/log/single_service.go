package log

import (
	"github.com/pkg/errors"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/app"
	"github.com/wxsatellite/goweb/framework/provider/config"
	"github.com/wxsatellite/goweb/framework/utils"
	"os"
	"path/filepath"
)

// SingleService 本地单个日志文件输出
type SingleService struct {
	BaseService

	folder string // 日志文件存储目录
	file   string // 日志文件名称
}

func NewSingleService(params ...interface{}) (interface{}, error) {
	if len(params) < 4 {
		return nil, errors.New("params error")
	}
	container := params[0].(framework.Container)
	level := params[1].(Level)
	ctxFielder := params[2].(CtxFielder)
	formatter := params[3].(Formatter)

	appService := container.MustMake(app.Key).(app.App)
	configService := container.MustMake(config.Key).(config.Config)

	service := &SingleService{}
	service.SetLevel(level)
	service.SetCtxFielder(ctxFielder)
	service.SetFormatter(formatter)

	// 获取默认的日志目录
	folder := appService.LogFolder()
	//  如果配置文件有配置，那么以配置文件为主
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}
	service.folder = folder

	// 不存在就创建
	if !utils.Exists(folder) {
		_ = os.MkdirAll(folder, os.ModePerm)
	}

	// 获取日志文件名称，以配置文件为主
	service.file = "goweb.log"
	if configService.IsExist("log.file") {
		service.file = configService.GetString("log.file")
	}

	// 打开日志文件句柄
	fd, err := os.OpenFile(filepath.Join(service.folder, service.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}

	// 设置输出管道
	service.SetOutput(fd)
	service.container = container
	return service, nil
}
