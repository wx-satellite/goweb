package framework

// NewInstance 服务容器通过这个函数创建服务提供者（ 这个函数由服务提供者提供 ）
type NewInstance func(...interface{}) (interface{}, error)

// ServiceProvider 服务提供者接口
type ServiceProvider interface {
	// Register 传入容器是因为后续如果希望根据一个服务的某个能力，比如配置服务的获取某个配置的能力，
	// 返回定义好的不同 NewInstance 函数，那我们就需要先从服务容器中获取配置服务，才能判断返回哪个 NewInstance。
	Register(Container) NewInstance

	// Boot 在调用实例化服务的时候会调用，可以把一些准备工作：基础配置，初始化参数的操作放在这个里面。
	// 如果 Boot 返回 error，整个服务实例化就会实例化失败，返回错误
	Boot(Container) error

	// IsDefer 决定是否在注册的时候实例化这个服务，如果不是注册的时候实例化，那就是在第一次 make 的时候进行实例化操作
	// false 表示不需要延迟实例化，在注册的时候就实例化，true 表示延迟实例化
	IsDefer() bool

	// Params 定义传递给 NewInstance 的参数，可以自定义多个，建议将 container 作为第一个参数
	Params(Container) []interface{}
	// Name 代表了这个服务提供者的凭证
	Name() string
}
