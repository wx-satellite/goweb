package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	postData := map[string]string{
		"email":    "18618192650",
		"from":     "https%3A%2F%2Fscrm.wxb.com%2F",
		"password": "asdf1234",
		"remember": "on",
	}
	urlParams := url.Values{}
	for key, value := range postData {
		urlParams.Set(key, value)
	}
	requestHandle, _ := http.NewRequest("POST", "https://account.wxb.com/index2/preLogin", strings.NewReader(urlParams.Encode()))
	requestHandle.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	requestHandle.Header.Set("Origin", "https://account.wxb.com")
	requestHandle.Header.Set("Content-Type", `application/x-www-form-urlencoded`)
	client := http.Client{}
	response, _ := client.Do(requestHandle)
	c := ""
	for _, cookie := range response.Cookies() {
		if cookie.Name == "PHPSESSID" {
			c += cookie.Name + "=" + cookie.Value + ";"
		}
	}
	requestHandle, _ = http.NewRequest("POST", "https://account.wxb.com/index2/login", strings.NewReader(urlParams.Encode()))
	requestHandle.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	requestHandle.Header.Set("Origin", "https://account.wxb.com")
	requestHandle.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requestHandle.Header.Set("Cookie", strings.Trim(c, ";"))
	response, _ = client.Do(requestHandle)
	requestHandle, _ = http.NewRequest("GET", "https://api-scrm.wxb.com/customer/list?page=2&page_size=20&corp_id=15708&corp_tag_ids=", nil)
	requestHandle.Header.Set("Cookie", strings.Trim(c, ";"))
	response, _ = client.Do(requestHandle)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.Status)
	fmt.Println(string(body))

	return
}
