package app

import (
	"errors"
	"github.com/google/uuid"
	flag "github.com/spf13/pflag"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/utils"
	"path/filepath"
)

type Service struct {
	container framework.Container // 服务容器

	baseFolder string // 基础路径：框架定义了业务目录最基本的几个路径，通过基础路径拼接可以获取到对应的目录路径

	appId string

	// 配置加载
	configMap map[string]string
}

func New(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)

	// 不存在就尝试从命令行参数中获取
	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base_folder", "", "项目的基础路径，默认为当前路径")
		flag.Parse()
	}

	// appId 为每一个应用的唯一标记，用于分布式锁
	return &Service{
		container:  container,
		baseFolder: baseFolder,
		appId:      uuid.New().String(),
		configMap:  make(map[string]string),
	}, nil
}

func (*Service) Version() string {
	return "1.0.0"
}

func (s *Service) BaseFolder() (path string) {
	if s.baseFolder != "" {
		return s.baseFolder
	}
	s.baseFolder = utils.GetExecDirectory()
	// 如果命令行参数没有值就获取当前路径
	return s.baseFolder
}

// ConfigFolder  表示配置文件地址
func (s *Service) ConfigFolder() string {
	if val, ok := s.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(s.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (s *Service) LogFolder() string {
	if val, ok := s.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(s.StorageFolder(), "log")
}

func (s *Service) HttpFolder() string {
	if val, ok := s.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(s.BaseFolder(), "http")
}

func (s *Service) ConsoleFolder() string {
	if val, ok := s.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(s.BaseFolder(), "console")
}

func (s *Service) StorageFolder() string {
	if val, ok := s.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(s.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (s *Service) ProviderFolder() string {
	if val, ok := s.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(s.BaseFolder(), "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (s *Service) MiddlewareFolder() string {
	if val, ok := s.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(s.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (s *Service) CommandFolder() string {
	if val, ok := s.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(s.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (s *Service) RuntimeFolder() string {
	if val, ok := s.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(s.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (s *Service) TestFolder() string {
	if val, ok := s.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(s.BaseFolder(), "test")
}

// AppId app的唯一标志
func (s *Service) AppId() string {
	return s.appId
}

// LoadAppConfig 加载配置map，用于更新app默认目录路径（从配置文件中获取）
func (s *Service) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		s.configMap[key] = val
	}
	return
}
