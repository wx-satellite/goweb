package log

import (
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/config"
	"github.com/wxsatellite/goweb/framework/provider/log/formatter"
	"io"
	"os"
	"strings"
)

type Provider struct {
	Driver string

	// 日志级别
	Level Level
	// 日志输出格式
	Formatter Formatter
	// 日志context上下文信息获取函数
	CtxFielder CtxFielder
	// 日志输出管道
	Output io.Writer
}

func (p *Provider) Register(container framework.Container) framework.NewInstance {
	if p.Driver == "" {
		configService, err := container.Make(config.Key)
		if err != nil {
			return NewConsoleService
		}
		service := configService.(config.Config)
		p.Driver = service.GetString("log.driver")
	}

	switch p.Driver {
	case "single":
		return NewSingleService
	case "rotate":
		return NewRotateService
	case "custom":
		return NewCustomService
	case "console":
		fallthrough
	default:
		return NewConsoleService
	}
}

func (p *Provider) Boot(container framework.Container) error {
	return nil
}

func (p *Provider) IsDefer() bool {
	return false
}

func (p *Provider) Name() string {
	return Key
}

func (p *Provider) Params(container framework.Container) []interface{} {

	// 获取config服务
	configService := container.MustMake(config.Key).(config.Config)

	if p.Level == 0 {
		p.Level = InfoLevel
		if configService.IsExist("log.level") {
			p.Level = logLevel(configService.GetString("log.level"))
		}
	}

	if p.Formatter == nil {
		p.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			switch configService.GetString("log.formatter") {
			case "json":
				p.Formatter = formatter.JsonFormatter
			case "text":
				fallthrough
			default:
				p.Formatter = formatter.TextFormatter
			}
		}
	}

	if p.CtxFielder == nil {
		p.CtxFielder = DefaultCtxFielder
	}

	if p.Output == nil {
		p.Output = os.Stdout
	}
	return []interface{}{container, p.Level, p.CtxFielder, p.Formatter, p.Output}
}

func logLevel(config string) Level {
	switch strings.ToLower(config) {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	case "trace":
		return TraceLevel
	}
	return InfoLevel
}
