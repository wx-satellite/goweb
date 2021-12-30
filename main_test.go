package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(filepath.Abs("1.txt"))
	return
}
