package distributed

import (
	"errors"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/app"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

/*
如果只有一个进程在跑定时任务或者一个机子在跑，那么如果这个机子故障或者进程异常退出会导致任务无法进行执行，容灾性低。
优化：一般会在一个机子上启动多个进程或者用多个机子分别跑进程，为了保证同一时间只有一个定时任务在运行因此需要引入分布式锁。
*/

// LocalService 分布式锁--文件锁：当一个服务器上有多个进程需要进行抢锁操作，文件锁是一种单机多进程抢占的很简易的实现方式
type LocalService struct {
	container framework.Container
}

func NewLocalProvider(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("param error")
	}
	container := params[0].(framework.Container)
	return &LocalService{container: container}, nil
}

func (s *LocalService) Select(serviceName string, appId string, holdTime time.Duration) (selectedAppId string, err error) {
	appService := s.container.MustMake(Key).(app.App)

	runtimeFolder := appService.RuntimeFolder()

	lockFile := filepath.Join(runtimeFolder, "distribute_"+serviceName)

	// 打开文件锁
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	// 尝试独占文件锁
	// LOCK_EX 放置互斥锁
	// LOCK_NB 非阻塞锁请求：默认情况下如果另一个进程持有了一把锁，那么调用syscall.Flock会被阻塞，如果设置为syscall.LOCK_NB就不会阻塞而是直接返回error。
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)

	// 获取不到锁，说明有进程先抢占了锁
	if err != nil {
		var selectAppIDByt []byte
		// 读取被选择的appId
		selectAppIDByt, err = ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByt), err
	}

	// 获取到了锁
	go func() {
		// 在一段时间之内，获取到锁是有效的，其他进程或者节点在这段时间内不能再进行抢占
		defer func() {
			_ = syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// 关闭文件
			_ = lock.Close()
			// 删除文件
			_ = os.Remove(lockFile)
		}()

		// 设置有效期：注意这里不使用sleep，gin框架的ShutDown也有类似的操作，如果使用sleep的话效率会比较低
		timer := time.NewTimer(holdTime)
		<-timer.C
	}()
	if _, err = lock.WriteString(appId); err != nil {
		return "", err
	}
	return appId, nil
}
