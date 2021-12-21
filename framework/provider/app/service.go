package app

import (
	"errors"
	"flag"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/utils"
	"path/filepath"
)

type Service struct {
	container framework.Container // 服务容器

	baseFolder string // 基础路径：框架定义了业务目录最基本的几个路径，通过基础路径拼接可以获取到对应的目录路径
}

func New(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	return &Service{container: container, baseFolder: baseFolder}, nil
}

func (*Service) Version() string {
	return "1.0.0"
}

func (s *Service) BaseFolder() (path string) {
	path = s.baseFolder
	if path != "" {
		return
	}

	// 不存在就尝试从命令行参数中获取
	var baseFolder string
	flag.StringVar(&baseFolder, "base_folder", "", "项目的基础路径，默认为当前路径")
	flag.Parse()

	path = baseFolder
	if path != "" {
		return
	}

	// 如果命令行参数没有值就获取当前路径
	return utils.GetExecDirectory()
}

// ConfigFolder  表示配置文件地址
func (s *Service) ConfigFolder() string {
	return filepath.Join(s.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (s *Service) LogFolder() string {
	return filepath.Join(s.StorageFolder(), "log")
}

func (s *Service) HttpFolder() string {
	return filepath.Join(s.BaseFolder(), "http")
}

func (s *Service) ConsoleFolder() string {
	return filepath.Join(s.BaseFolder(), "console")
}

func (s *Service) StorageFolder() string {
	return filepath.Join(s.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (s *Service) ProviderFolder() string {
	return filepath.Join(s.BaseFolder(), "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (s *Service) MiddlewareFolder() string {
	return filepath.Join(s.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (s *Service) CommandFolder() string {
	return filepath.Join(s.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (s *Service) RuntimeFolder() string {
	return filepath.Join(s.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (s *Service) TestFolder() string {
	return filepath.Join(s.BaseFolder(), "test")
}
