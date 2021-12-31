package log

import (
	"context"
	"io"
	"time"
)

/**
设计日志服务有三个思路：
	什么样的日志需要输出？ 对应日记级别
	日志输出哪些内容？ 对应日志格式
	日志输出到哪里？  比如控制台、文件等等
*/
const Key = "goweb:log"

type Level uint32

// 在 error 级别之上，我们把导致程序崩溃和导致请求结束的错误拆分出来，分为 panic 和 fatal 两个类型来定义级别
// 如果我们设置了日志输出级别为 info，那么 info级别以及info以上的级别日志也需要被打印出来
const (
	// UnknownLevel 表示未知的日志级别
	UnknownLevel Level = iota
	// PanicLevel level， panic 表示会导致整个程序出现崩溃的日志信息
	PanicLevel
	// FatalLevel level， fatal 表示会导致当前这个请求出现提前终止的错误信息
	FatalLevel
	// ErrorLevel level， error 表示出现错误，但是不一定影响后续请求逻辑的错误信息
	ErrorLevel
	// WarnLevel level， warn 表示出现错误，但是一定不影响后续请求逻辑的报警信息
	WarnLevel
	// InfoLevel level， info 表示正常的日志信息输出
	InfoLevel
	// DebugLevel level， debug 表示在调试状态下打印出来的日志信息
	DebugLevel
	//TraceLevel level， trace 表示最详细的信息，一般信息量比较大，可能包含调用堆栈等信息
	TraceLevel
)

// fields 表示当前日志的附带信息（上下文字段）
type Log interface {
	// Panic 表示会导致整个程序出现崩溃的日志信息
	Panic(ctx context.Context, msg string, fields map[string]interface{})
	// Fatal 表示会导致当前这个请求出现提前终止的错误信息
	Fatal(ctx context.Context, msg string, fields map[string]interface{})
	// Error 表示出现错误，但是不一定影响后续请求逻辑的错误信息
	Error(ctx context.Context, msg string, fields map[string]interface{})
	// Warn 表示出现错误，但是一定不影响后续请求逻辑的报警信息
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	// Info 表示正常的日志信息输出
	Info(ctx context.Context, msg string, fields map[string]interface{})
	// Debug 表示在调试状态下打印出来的日志信息
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	// Trace 表示最详细的信息，一般信息量比较大，可能包含调用堆栈等信息
	Trace(ctx context.Context, msg string, fields map[string]interface{})
	// SetLevel 设置日志级别
	SetLevel(level Level)

	// SetCtxFielder 从context中获取上下文字段field
	SetCtxFielder(handler CtxFielder)
	// SetFormatter 设置输出格式
	SetFormatter(formatter Formatter)

	// SetOutput 设置输出管道
	// 我们在定义接口的时候，并不知道它会输出到哪里。但是只需要知道一定会输出到某个输出管道（ 管道都会实现 io.Writer 接口 ）就可以了，
	// 之后在每个应用中使用的时候，我们再根据每个应用的配置，来确认具体的输出管道实现。
	SetOutput(out io.Writer)
}

/**
logger.Info(c, "demo test error", map[string]interface{}{
   "api":  "demo/demo",
   "user": "wx",
})
关于日志的上下文字段：它是一个 map 值，来源可能有两个：一个是用户在打印日志的时候传递的 map，比如上面代码中的 api 和 user；
而另外一部分数据是可能来自 context，因为在具体业务开发中，我们很有可能把一些通用信息，比如 trace_id 等放在 context 里，
这一部分信息也会希望取出放在日志的上下文字段中。
*/
// CtxFielder 定义了从context中获取信息的方法
type CtxFielder func(ctx context.Context) map[string]interface{}

func DefaultCtxFielder(ctx context.Context) map[string]interface{} {
	return nil
}

// Formatter 定义了将日志信息组织成字符串的通用方法
// eg："[Info]  2021-09-22T00:04:21+08:00    "demo test error"   map[api:demo/demo cspan_id: parent_id: span_id:c55051d94815vbl56i2g trace_id:c55051d94815vbl56i20 user:jianfengye]"
type Formatter func(level Level, t time.Time, msg string, fields map[string]interface{}) ([]byte, error)
