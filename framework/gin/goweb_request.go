package gin

import (
	"github.com/spf13/cast"
)

/**
让 gin 的 Context 实现 IRequest 和 IResponse 接口，可以发现接口存在以下三种情况：
	Gin 框架中已经存在相同参数、相同返回值、相同接口名的接口
	Gin 框架中不存在相同接口名的接口
	Gin 框架中存在相同接口名，但是不同返回值或者不同参数的接口

对于第三种情况则棘手一些。以 Gin 中已经有的 QueryXXX 系列接口为例，它的 QueryXXX 系列接口和我们想要的有一定差别，
比如它的 Query 不支持多种数据类型的直接获取。解决方案：将 IRequest 的 QueryXXX 系列接口重新命名，
又因为接口都带有一个默认值，所以我们将其重新命名为 DefaultQueryXXX。
*/

// const defaultMultipartMemory = 32 << 20 // 32 MB

// 代表请求包含的方法
type IRequest interface {
	// 请求地址url中带的参数
	// 形如: foo.com?a=1&b=bar&c[]=bar
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat64(key string, def float64) (float64, bool)
	DefaultQueryFloat32(key string, def float32) (float32, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStringSlice(key string, def []string) ([]string, bool)

	// 路由匹配中带的参数
	// 形如 /book/:id
	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat64(key string, def float64) (float64, bool)
	DefaultParamFloat32(key string, def float32) (float32, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	DefaultParamString(key string, def string) (string, bool)
	DefaultParam(key string) interface{}

	// form表单中带的参数
	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat64(key string, def float64) (float64, bool)
	DefaultFormFloat32(key string, def float32) (float32, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormString(key string, def string) (string, bool)
	DefaultFormStringSlice(key string, def []string) ([]string, bool)
	//DefaultFormFile(key string) (*multipart.FileHeader, error)
	DefaultForm(key string) interface{}

	// json body
	//BindJson(obj interface{}) error
	//
	//// xml body
	//BindXml(obj interface{}) error
	//
	//// 其他格式
	//GetRawData() ([]byte, error)
	//
	//// 基础信息
	//Uri() string
	//Method() string
	//Host() string
	//ClientIp() string
	//
	//// header
	//Headers() map[string]string
	//Header(key string) (string, bool)
	//
	//// cookie
	//Cookies() map[string]string
	//Cookie(key string) (string, bool)
}

/**
Gin 框架在获取 QueryAll 的时候，使用了一个 QueryCache，实现了在第一次获取参数的时候，
用方法 initQueryCache 将 ctx.request.URL.Query() 缓存在内存中：代码实现如下，比较简单
func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Request != nil {
			c.queryCache = c.Request.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}
*/
// 获取请求地址中所有参数
func (c *Context) QueryAll() map[string][]string {
	c.initQueryCache()
	return map[string][]string(c.queryCache)
}

// 请求地址url中带的参数
// 形如: foo.com?a=1&b=bar&c[]=bar

// 获取Int类型的请求参数
func (c *Context) DefaultQueryInt(key string, def int) (int, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			// 使用cast库将string转换为Int
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryFloat64(key string, def float64) (float64, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryFloat32(key string, def float32) (float32, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryString(key string, def string) (string, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return def, false
}

func (c *Context) DefaultQueryStringSlice(key string, def []string) ([]string, bool) {
	params := c.QueryAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

// 路由匹配中带的参数
// 形如 /book/:id
func (c *Context) DefaultParamInt(key string, def int) (int, bool) {
	if val := c.DefaultParam(key); val != nil {
		// 通过cast进行类型转换
		return cast.ToInt(val), true
	}
	return def, false
}

func (c *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	if val := c.DefaultParam(key); val != nil {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (c *Context) DefaultParamFloat64(key string, def float64) (float64, bool) {
	if val := c.DefaultParam(key); val != nil {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (c *Context) DefaultParamFloat32(key string, def float32) (float32, bool) {
	if val := c.DefaultParam(key); val != nil {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (c *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	if val := c.DefaultParam(key); val != nil {
		return cast.ToBool(val), true
	}
	return def, false
}

func (c *Context) DefaultParamString(key string, def string) (string, bool) {
	if val := c.DefaultParam(key); val != nil {
		return cast.ToString(val), true
	}
	return def, false
}

// 获取路由参数
func (c *Context) DefaultParam(key string) interface{} {
	if val, ok := c.Params.Get(key); ok {
		return val
	}
	return nil
}

func (c *Context) FormAll() map[string][]string {
	c.initFormCache()
	return map[string][]string(c.formCache)
}

func (c *Context) DefaultFormInt(key string, def int) (int, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultFormFloat64(key string, def float64) (float64, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultFormFloat32(key string, def float32) (float32, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (c *Context) DefaultFormStringSlice(key string, def []string) ([]string, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

func (c *Context) DefaultFormString(key string, def string) (string, bool) {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return def, false
}

func (c *Context) DefaultForm(key string) interface{} {
	params := c.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return nil
}
