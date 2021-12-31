package formatter

import "github.com/wxsatellite/goweb/framework/provider/log"

func Prefix(level log.Level) (prefix string) {
	switch level {
	case log.PanicLevel:
		prefix = "[Panic]"
	case log.FatalLevel:
		prefix = "[Fatal]"
	case log.ErrorLevel:
		prefix = "[Error]"
	case log.WarnLevel:
		prefix = "[Warn]"
	case log.InfoLevel:
		prefix = "[Info]"
	case log.DebugLevel:
		prefix = "[Debug]"
	case log.TraceLevel:
		prefix = "[Trace]"
	}
	return
}
