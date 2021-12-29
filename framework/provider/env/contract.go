package env

const (
	// Production 生产环境
	Production = "production"

	// Testing 生产环境
	Testing = "testing"

	// Development 开发环境
	Development = "development"

	Key = "goweb:env"
)

/**
环境变量可能会有很多。但是我们每次部署一个环境的时候，设置的环境变量可能就只有一两个，那其他的环境变量就需要有一个“默认值”。
这个默认值我们一般使用一个以 dot 点号开头的文件.env 来进行设置

环境变量是当前程序运行时的环境，一个程序运行时就是固定的，因此不允许在运行中修改，也就没有了 Set 方法。


*/
type Env interface {
	// AppEnv 获取当前的环境，建议分为 development/testing/production
	AppEnv() string
	// IsExist 判断一个环境变量是否有被设置
	IsExist(key string) bool
	// Get 获取某个环境变量，如果没有设置，返回""
	Get(string) string
	// All 获取所有的环境变量，.env 和运行环境变量融合后结果
	All() map[string]string
}
