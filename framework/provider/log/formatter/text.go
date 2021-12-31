package formatter

import (
	"bytes"
	"fmt"
	"github.com/wxsatellite/goweb/framework/provider/log"
	"time"
)

func TextFormatter(level log.Level, t time.Time, msg string, fields map[string]interface{}) (res []byte, err error) {
	bf := bytes.NewBuffer([]byte{})

	separator := "\t"

	bf.WriteString(Prefix(level))
	bf.WriteString(separator)

	// 输出时间
	bf.WriteString(t.Format(time.RFC3339))
	bf.WriteString(separator)

	// 输出文本
	bf.WriteString("\"")
	bf.WriteString(msg)
	bf.WriteString("\"")
	bf.WriteString(separator)

	// 输出map
	// fmt.Sprint() 输出：map[name:wx]
	bf.WriteString(fmt.Sprint(fields))
	return bf.Bytes(), nil
}
