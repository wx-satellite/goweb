package config

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/app"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// TODO: 实现可以通过类似这样子的命令 `./goweb config get "database.mysql"` 获取配置信息

// Service 配置文件服务
// database.mysql.password 获取database.yaml文件中mysql配置对应的password字段
type Service struct {
	container framework.Container
	folder    string // 配置文件目录
	keyBreak  string // 路径分隔符，默认是 "."

	envMaps  map[string]string      // 所有环境变量
	confMaps map[string]interface{} // 配置文件结构，key为文件名
	confRaws map[string][]byte      // 配置文件的原始信息

	// 由于在运行时增加了对 confMaps 的写操作（配置文件热更新）所以需要对 confMaps 进行锁设置，以防止在写 confMaps 的时候，读操作进入读取了错误信息。
	// 其次：目前这个场景，读明显多于写。所以我们的锁是一个读写锁，读写锁可以让多个读并发读，但是只要有一个写操作，读和写都需要等待。
	lock sync.RWMutex
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
		// key是文件名，value是yaml.Unmarshal的结果
		confMaps: make(map[string]interface{}),
		// key是文件名，value是文件内容
		confRaws: make(map[string][]byte),
		lock:     sync.RWMutex{},
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
		if err = service.handleConfigFile(file.Name(), envFolder); err != nil {
			fmt.Println(fmt.Sprintf("加载配置文件失败：%s，错误：%v", file.Name(), err))
			continue
		}
	}

	/* 监控配置文件的修改 */
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	// 监控配置文件目录下文件的修改
	err = watch.Add(envFolder)
	if err != nil {
		// fsnotify使用了操作系统接口，监听器中保存了系统资源的句柄，所以使用后需要关闭
		_ = watch.Close()
		return nil, err
	}
	go func() {
		defer func() {
			_ = watch.Close()
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				// 判断事件的类型
				// ev.Name 的值类似：/Users/weixin/Desktop/goweb/config/development/app.yml
				path, _ := filepath.Abs(ev.Name) // 获取绝对路径，如果已经是绝对路径就不做任何操作
				index := strings.LastIndex(path, string(os.PathSeparator))
				folder := path[:index]
				fileName := path[index+1:]
				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("创建文件 : ", ev.Name)
					_ = service.handleConfigFile(fileName, folder)
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", ev.Name)
					_ = service.handleConfigFile(fileName, folder)
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("删除文件 : ", ev.Name)
					_ = service.removeConfigFile(fileName, folder)
				}
			case err := <-watch.Errors:
				log.Println("监控配置文件错误：", err)
			}
		}
	}()
	return service, nil
}

// removeConfigFile 删除内存中的配置文件信息
func (s *Service) removeConfigFile(fileName string, envFolder string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	splits := strings.Split(fileName, ".")
	// 只处理 yaml 或者 yml 结尾的配置文件
	if !(len(splits) == 2 && (splits[1] == "yaml" || splits[1] == "yml")) {
		return errors.New("只支持处理yaml、yml后缀的文件")
	}
	// 删除
	delete(s.confMaps, splits[0])
	delete(s.confRaws, splits[0])
	return
}

// handleConfigFile 更新内存中的配置文件信息
func (s *Service) handleConfigFile(fileName string, envFolder string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	splits := strings.Split(fileName, ".")
	// 只处理 yaml 或者 yml 结尾的配置文件
	if !(len(splits) == 2 && (splits[1] == "yaml" || splits[1] == "yml")) {
		return errors.New("只支持处理yaml、yml后缀的文件")
	}
	name := splits[0]
	// 读取文件内容
	var bf []byte
	bf, err = ioutil.ReadFile(filepath.Join(envFolder, fileName))
	if err != nil {
		return
	}
	s.confRaws[name] = bf

	// 将环境变量占位替换成环境变量的值
	bf = replace(bf, s.envMaps)

	// 解析 yaml
	c := make(map[string]interface{})

	if err = yaml.Unmarshal(bf, &c); err != nil {
		return
	}
	// 文件名为key
	s.confMaps[name] = c

	// 如果文件是 app.yml 那么需要更新一下app服务的默认目录路径
	if name == "app" && s.container.IsBind(app.Key) {
		if path, ok := c["path"]; ok {
			appService := s.container.MustMake(app.Key).(app.App)
			appService.LoadAppConfig(cast.ToStringMapString(path))
		}
	}
	return
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
	s.lock.RLock()
	defer s.lock.RUnlock()
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
