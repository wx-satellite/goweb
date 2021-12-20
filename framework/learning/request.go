package learning

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"io/ioutil"
	"mime/multipart"
)

/*
对比我们自己的query方法和gin的query差别：
	gin的query方法，通过本地内存cache缓存了请求的query参数。后续每次读query参数时，都会从内存中直接读，
减少了每次都要调用request.Query()方法的开销。

gin 有很多实现细节值得我们学习和拜读！
*/
const defaultMultipartMemory = 32 << 20 // 32MB
// IRequest 将 request 相关的功能声明成接口
// 关于接口：对于比较完整的功能模块，先定义接口再具体实现。一个清晰的接口可以让使用者非常清楚，这个功能
// 提供了哪些函数，哪些函数是我需要的，用的时候方便，查找的时候也方便。其次接口可以实现解耦，使用者在用的时候没有心里负担，即使 "实现改变了也不需要修改代码"。
type IRequest interface {

	// 参数有两个：
	//	一个是 key，代表从参数列表中查找指定的key
	// 	另外一个是 def，代表如果查找不到则使用默认值返回。
	// 返回值有两个：
	//	一个代表对应 key 的匹配值
	//	另一个 bool 代表是否匹配

	// 解析请求url中携带的参数
	// 例如：foo.com?a=1&b=2
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)

	// 匹配路由中的参数
	// 例如：/book/:id
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)

	// form 表单的参数
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)

	// json body
	BindJson(obj interface{}) error

	// xml body
	BindXml(obj interface{}) error

	// 其他格式
	GetRawData() ([]byte, error)

	// 基础信息
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	// header
	Headers() map[string][]string
	Header(key string) (string, bool)

	// cookie
	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

// QueryAll 获取请求地址中的所有参数
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}
	return ctx.request.URL.Query()
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToInt(val[0]), true
}

func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToInt64(val[0]), true
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToFloat32(val[0]), true
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToFloat64(val[0]), true
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToBool(val[0]), true
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return val[0], true
}

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.QueryAll()
	val, ok := params[key]
	if !ok {
		return def, false
	}
	return val, true
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request == nil {
		return nil
	}
	// Form：存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。
	// PostForm：存储了 post、put 参数，在使用之前需要调用 ParseForm 方法。
	_ = ctx.request.ParseForm()
	return ctx.request.PostForm
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	forms := ctx.FormAll()
	val, ok := forms[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToInt64(val[0]), true
}

func (ctx *Context) FormInt(key string, def int) (int, bool) {
	forms := ctx.FormAll()
	val, ok := forms[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToInt(val[0]), true
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	forms := ctx.FormAll()
	val, ok := forms[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToFloat64(val[0]), true
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	forms := ctx.FormAll()
	val, ok := forms[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToFloat32(val[0]), true
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	forms := ctx.FormAll()
	val, ok := forms[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return cast.ToBool(val[0]), true
}

func (ctx *Context) FormString(key string, def string) (string, bool) {
	forms := ctx.FormAll()
	val, ok := forms[key]
	if !ok || len(val) <= 0 {
		return def, false
	}
	return val[0], true
}

func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	forms := ctx.FormAll()
	val, ok := forms[key]
	if !ok {
		return def, false
	}
	return val, true
}

func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.request.MultipartForm == nil {
		// 使用ParseMultipartForm()方法解析Form，解析时会读取所有数据，但需要指定保存在内存中的最大字节数，剩余的字节数会保存在临时磁盘文件中。
		if err := ctx.request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	// FormFile 的实现当 MultipartForm 为nil的时候会自动调用 ParseMultipartForm 方法
	f, fh, err := ctx.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	//  multipart.FileHeader 如何使用可以查看 FormFile 的源码 ---> fhs[0].Open()
	_ = f.Close()
	return fh, nil
}

// BindJson 将body解析到obj结构体中
func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request == nil {
		return errors.New("ctx.request is empty")
	}
	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return err
	}

	// request.Body 的读取是一次性的，读取一次之后，下个逻辑再去 request.Body 中是读取不到数据内容的。
	// 所以我们读取完 request.Body 之后，还要再复制一份 Body 内容，填充到 request.Body 里。
	// 重新填充request.Body，为后续的逻辑二次读取做准备
	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// 解析到结构体中
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) BindXml(obj interface{}) error {
	if ctx.request == nil {
		return errors.New("ctx.request is empty")
	}
	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return err
	}
	// 重新填充request.Body，为后续的逻辑二次读取做准备
	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// 解析到结构体中
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *Context) GetRawData() ([]byte, error) {
	if ctx.request == nil {
		return nil, errors.New("ctx.request is empty")
	}
	body, err := ioutil.ReadAll(ctx.request.Body)
	if err != nil {
		return nil, err
	}
	// 重新填充request.Body，为后续的逻辑二次读取做准备
	ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// Uri 会携带get参数，例如： /name?age=12&hobby=basketball
func (ctx *Context) Uri() string {
	if ctx.request == nil {
		return ""
	}
	return ctx.request.RequestURI
}

func (ctx *Context) Method() string {
	if ctx.request == nil {
		return ""
	}
	return ctx.request.Method
}

func (ctx *Context) Host() string {
	if ctx.request == nil {
		return ""
	}
	return ctx.request.Host
}

func (ctx *Context) ClientIp() string {
	if ctx.request == nil {
		return ""
	}
	ipAddress := ctx.request.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = ctx.request.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = ctx.request.RemoteAddr
	}
	return ipAddress
}

func (ctx *Context) Headers() map[string][]string {
	if ctx.request == nil {
		return nil
	}
	return ctx.request.Header
}

func (ctx *Context) Header(key string) (string, bool) {
	if ctx.request == nil {
		return "", false
	}
	val := ctx.request.Header.Values(key)
	if len(val) <= 0 {
		return "", false
	}
	return val[0], true
}

func (ctx *Context) Cookies() map[string]string {
	if ctx.request == nil {
		return nil
	}
	cookies := ctx.request.Cookies()
	res := make(map[string]string)
	for _, cookie := range cookies {
		res[cookie.Name] = cookie.Value
	}
	return res
}

func (ctx *Context) Cookie(key string) (string, bool) {
	cookies := ctx.Cookies()
	if len(cookies) <= 0 {
		return "", false
	}
	val, ok := cookies[key]
	return val, ok
}

func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	val, ok := ctx.params[key]
	if !ok {
		return def, false
	}
	return cast.ToInt(val), true
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	val, ok := ctx.params[key]
	if !ok {
		return def, false
	}
	return cast.ToInt64(val), true
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	val, ok := ctx.params[key]
	if !ok {
		return def, false
	}
	return cast.ToFloat32(val), true
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	val, ok := ctx.params[key]
	if !ok {
		return def, false
	}
	return cast.ToFloat64(val), true
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	val, ok := ctx.params[key]
	if !ok {
		return def, false
	}
	return cast.ToBool(val), true
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	val, ok := ctx.params[key]
	if !ok {
		return def, false
	}
	return val, true
}
