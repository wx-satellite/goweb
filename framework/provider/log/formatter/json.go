package formatter

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/wxsatellite/goweb/framework/provider/log"
	"time"
)

func JsonFormatter(level log.Level, t time.Time, msg string, fields map[string]interface{}) (res []byte, err error) {
	bf := bytes.NewBuffer([]byte{})

	fields["level"] = level
	fields["msg"] = msg
	fields["timestamp"] = t.Format(time.RFC3339)
	c, err := json.Marshal(fields)
	if err != nil {
		return bf.Bytes(), errors.Wrap(err, "json format error")
	}
	bf.Write(c)
	return bf.Bytes(), nil
}
