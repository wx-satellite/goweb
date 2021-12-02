package main

import (
	"fmt"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	uri := "name"
	fmt.Println(strings.SplitN(uri, "/", 2))
}
