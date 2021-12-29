package env

import (
	"bufio"
	"errors"
	"github.com/wxsatellite/goweb/framework"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	folder    string // .env所在的目录
	container framework.Container
	envs      map[string]string // 保存所有的环境变量
}

func New(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	folder := params[1].(string)
	server := &Service{
		container: container,
		folder:    folder,
		// APP_ENV 参考 vue 和 laravel 都会设置这么一个固定的环境变量，同时预设了三个模式：开发、测试、生产
		envs: map[string]string{"APP_ENV": Development},
	}

	// 解析 .env 文件
	file := filepath.Join(folder, ".env")

	// 打开 .env 文件，并获取默认环境变量
	fd, err := os.Open(file)
	if err == nil {
		defer func() {
			_ = fd.Close()
		}()
		br := bufio.NewReader(fd)
		for {
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			// APP_ENV=2
			res := strings.Split(string(line), "=")
			// 不符合规范，过滤
			if len(res) < 2 {
				continue
			}
			server.envs[strings.TrimSpace(res[0])] = strings.TrimSpace(res[1])
		}
	}

	// 获取运行时设置的环境变量，会覆盖 .env
	for _, e := range os.Environ() {
		res := strings.Split(e, "=")
		// 不符合规范，过滤
		if len(res) < 2 {
			continue
		}
		server.envs[strings.TrimSpace(res[0])] = strings.TrimSpace(res[1])
	}
	return server, nil
}

func (s *Service) AppEnv() string {
	return s.Get("APP_ENV")
}

func (s *Service) Get(key string) (value string) {
	value = s.envs[key]
	return
}
func (s *Service) All() map[string]string {
	return s.envs
}
func (s *Service) IsExist(key string) (ok bool) {
	_, ok = s.envs[key]
	return
}
