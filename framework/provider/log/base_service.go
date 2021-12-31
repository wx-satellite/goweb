package log

import (
	"context"
	"github.com/wxsatellite/goweb/framework"
	"github.com/wxsatellite/goweb/framework/provider/log/formatter"
	"io"
	"log"
	"time"
)

type BaseService struct {
	level      Level               // 日志级别
	formatter  Formatter           // 日志格式化方法
	ctxFielder CtxFielder          // ctx获取上下文字段
	output     io.Writer           // 输出
	container  framework.Container // 容器
}

// level 越小级别越大
func (s *BaseService) CanLog(level Level) bool {
	return s.level >= level
}

// logf 先判断日志级别是否符合要求，如果不符合要求，则直接返回，不进行打印；
// 再使用 ctxFielder，从 context 中获取信息放在上下文字段中；
// 接着将日志信息按照 formatter 序列化为字符串；
// 最后通过 output 进行输出。
func (s *BaseService) logf(level Level, ctx context.Context, msg string, fields map[string]interface{}) (err error) {
	if !s.CanLog(level) {
		return
	}

	// 用户传入的上下文字段
	currentFields := fields

	// 填充 context 中的一些通用上下文字段
	if s.ctxFielder != nil {
		values := s.ctxFielder(ctx)
		for key, value := range values {
			currentFields[key] = value
		}
	}

	// 如果没有格式化函数，那么默认以文本的形式格式化
	if s.formatter == nil {
		s.formatter = formatter.TextFormatter
	}

	// 序列化日志信息
	bs, err := s.formatter(level, time.Now(), msg, currentFields)
	if err != nil {
		return
	}

	// 如果是 panic 级别就使用标准库的log.Panicln
	if level == PanicLevel {
		log.Panicln(string(bs))
		return nil
	}

	// 通过 output 进行输出
	_, _ = s.output.Write(bs)
	_, _ = s.output.Write([]byte("\r\n"))
	return
}

// Panic 输出panic的日志信息
func (s *BaseService) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = s.logf(PanicLevel, ctx, msg, fields)
}

// Fatal will add fatal record which contains msg and fields
func (s *BaseService) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = s.logf(FatalLevel, ctx, msg, fields)
}

// Error will add error record which contains msg and fields
func (s *BaseService) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = s.logf(ErrorLevel, ctx, msg, fields)
}

// Warn will add warn record which contains msg and fields
func (s *BaseService) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = s.logf(WarnLevel, ctx, msg, fields)
}

// Info 会打印出普通的日志信息
func (s *BaseService) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = s.logf(InfoLevel, ctx, msg, fields)
}

// Debug will add debug record which contains msg and fields
func (s *BaseService) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = s.logf(DebugLevel, ctx, msg, fields)
}

// Trace will add trace info which contains msg and fields
func (s *BaseService) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = s.logf(TraceLevel, ctx, msg, fields)
}

// SetOutput 设置输出管道
func (s *BaseService) SetOutput(output io.Writer) {
	s.output = output
}

// SetLevel 设置日志级别
func (s *BaseService) SetLevel(level Level) {
	s.level = level
}

// SetFormatter 设置序列化函数
func (s *BaseService) SetFormatter(formatter Formatter) {
	s.formatter = formatter
}

// SetCtxFielder 设置从context获取字段的函数
func (s *BaseService) SetCtxFielder(handler CtxFielder) {
	s.ctxFielder = handler
}
