package framework

import (
	"errors"
	"sync"
)

type Container interface {
	// Bind 绑定服务提供者，如果已经存在则会替换
	Bind(provider ServiceProvider) error

	// IsBind 判断是否绑定过指定的服务提供者
	IsBind(key string) bool

	// Make 根据关键字获取一个服务
	Make(key string) (interface{}, error)

	// MustMake 获取服务实例，如果服务没有绑定则会panic
	MustMake(key string) interface{}

	// MakeNew 获取服务实例，只是这个服务并不是单例模式的，它是根据服务提供者和传递的 params 参数实例化出来的
	MakeNew(key string, params []interface{}) (interface{}, error)
}

type GoWebContainer struct {
	Container

	// providers 存储注册的服务提供者，key 为字符串凭证
	providers map[string]ServiceProvider
	// instance 存储具体的实例，key 为字符串凭证
	instances map[string]interface{}
	// lock 容器的场景是读多写少的，因此使用读些锁而不是互斥锁
	lock sync.RWMutex
}

func NewGoWebContainer() *GoWebContainer {
	return &GoWebContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]interface{}),
		lock:      sync.RWMutex{},
	}
}

// Bind 绑定服务提供者
func (c *GoWebContainer) Bind(provider ServiceProvider) (err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	key := provider.Name()

	c.providers[key] = provider

	// 延迟实例化就直接返回
	if provider.IsDefer() {
		return
	}
	instance, err := c.newInstance(provider, nil)
	if err != nil {
		return
	}
	c.instances[key] = instance
	return
}

// Make 创建服务，会缓存
func (c *GoWebContainer) Make(key string) (res interface{}, err error) {
	return c.make(key, nil, false)
}

// MustMake 创建服务，会缓存，但是出错会panic
func (c *GoWebContainer) MustMake(key string) (res interface{}) {
	res, err := c.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return
}

// MakeNew 创建指定参数的服务，不走缓存
func (c *GoWebContainer) MakeNew(key string, params []interface{}) (res interface{}, err error) {
	return c.make(key, params, true)
}

// IsBind 是否已经绑定
func (c *GoWebContainer) IsBind(key string) bool {
	return c.findServiceProvider(key) != nil
}

// findServiceProvider 获取服务提供者
func (c *GoWebContainer) findServiceProvider(key string) (sp ServiceProvider) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	sp, _ = c.providers[key]
	return
}

// make 实例化操作
func (c *GoWebContainer) make(key string, params []interface{}, forceNew bool) (res interface{}, err error) {
	c.lock.RLock()

	// 不强制重新初始化，就获取容器中存在的实例
	res, ok := c.instances[key]
	if ok {
		c.lock.RUnlock()
		return
	}

	// 获取服务提供者
	provider := c.findServiceProvider(key)
	if provider == nil {
		c.lock.RUnlock()
		err = errors.New("contract " + key + " have not register")
		return
	}
	// 关闭读锁
	c.lock.RUnlock()

	// 加写锁，双重检查，在获得写锁之后，可能有别的协程已经创建完成，可以直接返回，避免再创建
	c.lock.Lock()
	defer c.lock.Unlock()
	res, ok = c.instances[key]
	if ok {
		return
	}

	if forceNew {
		return c.newInstance(provider, params)
	}

	// 如果容器不存在就新增一个
	newInstance, err := c.newInstance(provider, params)
	if err != nil {
		return
	}
	c.instances[key] = newInstance
	res = newInstance
	return
}

// newInstance 实例化操作
func (c *GoWebContainer) newInstance(sp ServiceProvider, params []interface{}) (res interface{}, err error) {
	if err = sp.Boot(c); err != nil {
		return
	}

	// 如果没有指定的实例化参数，就获取sp自带的
	if params == nil {
		params = sp.Params(c)
	}

	// 获取实例化方法
	newInstanceMethod := sp.Register(c)

	// 实例化
	res, err = newInstanceMethod(params...)
	return
}
