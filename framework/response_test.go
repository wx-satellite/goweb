package framework

import (
	"fmt"
	"strings"
	"testing"
)

func TestContext_Html(t *testing.T) {
	file := "name/article.tmpl"
	paths := strings.Split(file, "/")
	if len(paths) <= 0 {
		fmt.Println("paths is empty")
		return
	}
	lastPath := paths[len(paths)-1]
	name := lastPath[:strings.Index(lastPath, ".")]
	fmt.Println(name)
}
