package learning

import (
	"testing"
)

func Test(t *testing.T) {
	core := NewCore()
	group := core.Group("/user")
	group.Get("/name", nil)
	group.Group("/haha").Get("/name", nil)
}
