package config

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/wxsatellite/goweb/framework"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Service 配置文件服务
// database.mysql.password 获取database.yaml文件中mysql配置对应的password字段
type Service struct {
	container framework.Container
	folder    string // 配置文件目录
	keyBreak  string // 路径分隔符，默认是 "."

	envMaps  map[string]string      // 所有环境变量
	confMaps map[string]interface{} // 配置文件结构，key为文件名
	confRaws map[string][]byte      // 配置文件的原始信息
}

func New(params ...interface{}) (interface{}, error) {
	if len(params) < 4 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	folder := params[1].(string)
	env := params[2].(string)
	envMaps := params[3].(map[string]string)

	envFolder := filepath.Join(folder, env)

	// 检测配置文件路径是否存在
	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return nil, errors.New("folder " + envFolder + " not exist: " + err.Error())
	}
	service := &Service{
		container: container,
		folder:    folder,
		keyBreak:  ".",
		envMaps:   envMaps,
		confMaps:  make(map[string]interface{}),
		confRaws:  make(map[string][]byte),
	}

	//  获取目录下所有的配置文件
	files, err := os.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 获取每一个文件
	for _, file := range files {
		// 只读取文件
		if file.IsDir() {
			continue
		}
		splits := strings.Split(file.Name(), ".")
		// 只处理 yaml 或者 yml 结尾的配置文件
		if !(len(splits) == 2 && (splits[1] == "yaml" || splits[1] == "yml")) {
			continue
		}
		name := splits[0]
		// 读取文件内容
		var bf []byte
		bf, err = ioutil.ReadFile(filepath.Join(envFolder, name))
		if err != nil {
			continue
		}
		service.confRaws[name] = bf

		// 将环境变量占位替换成环境变量的值
		bf = replace(bf, envMaps)

		// 解析 yaml
		c := make(map[string]interface{})

		if err = yaml.Unmarshal(bf, &c); err != nil {
			continue
		}
		// 文件名为key
		service.confMaps[name] = c
	}
	return service, nil
}

// IsExist check setting is exist
func (s *Service) IsExist(key string) bool {
	return s.find(key) != nil
}

// Get a new interface
func (s *Service) Get(key string) interface{} {
	return s.find(key)
}

// GetBool get bool type
func (s *Service) GetBool(key string) bool {
	return cast.ToBool(s.find(key))
}

// GetInt get Int type
func (s *Service) GetInt(key string) int {
	return cast.ToInt(s.find(key))
}

// GetFloat64 get float64
func (s *Service) GetFloat64(key string) float64 {
	return cast.ToFloat64(s.find(key))
}

// GetTime get time type
func (s *Service) GetTime(key string) time.Time {
	return cast.ToTime(s.find(key))
}

// GetString get string typen
func (s *Service) GetString(key string) string {
	return cast.ToString(s.find(key))
}

// GetIntSlice get int slice type
func (s *Service) GetIntSlice(key string) []int {
	return cast.ToIntSlice(s.find(key))
}

// GetStringSlice get string slice type
func (s *Service) GetStringSlice(key string) []string {
	return cast.ToStringSlice(s.find(key))
}

// GetStringMap get map which key is string, value is interface
func (s *Service) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(s.find(key))
}

// GetStringMapString get map which key is string, value is string
func (s *Service) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(s.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (s *Service) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(s.find(key))
}

// Load a config to a struct, val should be an pointer
func (s *Service) Load(key string, val interface{}) error {
	return mapstructure.Decode(s.find(key), val)
}

// replace 配置文件也会使用环境变量的值，使用：env(xxx) 占位，因此解析配置文件的时候，需要替换成实际的环境变量值
func replace(content []byte, envMaps map[string]string) []byte {
	if envMaps == nil {
		return content
	}
	for key, val := range envMaps {
		reKey := fmt.Sprintf("env(%s)", key)
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}
	return content
}

// find 获取某一个配置
func (s *Service) find(key string) interface{} {
	return searchMap(s.confMaps, strings.Split(key, s.keyBreak))
}

// searchMap 从配置文件中获取指定的配置
// path = [database mysql password]表示database.yaml中获取mysql配置项的password字段
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}
	configItem, ok := source[path[0]]
	if !ok {
		return nil
	}
	// 当path只有1时，就直接返回
	if len(path) == 1 {
		return configItem
	}

	// 继续递归查找，直到 len(params)=1 就返回
	switch v := configItem.(type) {
	case map[interface{}]interface{}:
		return searchMap(cast.ToStringMap(v), path[1:])
	case map[string]interface{}:
		return searchMap(v, path[1:])
	default:
		return nil
	}
}
