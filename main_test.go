package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	bs, _ := json.Marshal(map[string]string{"name": "bob", "age": "12"})
	request, _ := http.NewRequest("POST", "http://www.baidu.com", bytes.NewBuffer(bs))
	httpRequest, _ := httputil.DumpRequest(request, true)
	fmt.Println(string(httpRequest))
	headers := strings.Split(string(httpRequest), "\r\n")
	fmt.Println(headers)
}
