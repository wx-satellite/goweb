package app

// 应用的目录服务：框架规范业务目录虽然看起来有点越俎代庖，但是规范某些最小化的功能性目录还是很有必要的

const Key = "goweb:app"

// App 定义接口
// 根据 App 的接口定义我们可以得出一个业务目录应该需要如下目录：配置文件、日志、服务提供者与实现、中间件、命令行、运行时产生的信息、测试等目录。
type App interface {
	// Version 定义当前版本
	Version() string
	//BaseFolder 定义项目基础地址
	BaseFolder() string
	// ConfigFolder 定义了配置文件的路径
	ConfigFolder() string
	// LogFolder 定义了日志所在路径
	LogFolder() string
	// ProviderFolder 定义业务自己的服务提供者地址
	ProviderFolder() string
	// MiddlewareFolder 定义业务自己定义的中间件
	MiddlewareFolder() string
	// CommandFolder 定义业务定义的命令
	CommandFolder() string
	// RuntimeFolder 定义业务的运行中间态信息
	RuntimeFolder() string
	// TestFolder 存放测试所需要的信息
	TestFolder() string

	// AppId 表示当前app的唯一id，用于分布式锁
	AppId() string

	// LoadAppConfig 根据配置文件更新默认目录路径
	LoadAppConfig(kv map[string]string)
}
