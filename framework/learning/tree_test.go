package learning

import (
	"fmt"
	"strings"
	"testing"
)

func TestTree_AddRouter(t *testing.T) {
	tree := NewTree()
	fmt.Println(tree.AddRouter("/:user/name", nil))
	fmt.Println(tree.AddRouter("/:user/name/:age", nil))
}

func TestSplit(t *testing.T) {
	fmt.Println(strings.SplitN(":user/name", "/", 2))
	fmt.Println(strings.SplitN("/:user/name", "/", 2))
}
