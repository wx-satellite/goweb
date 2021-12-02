package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s := make([]int64, 0, 4)
	s = append(s, 1, 2)

	s1 := append(s, 3)

	fmt.Println(s)
	fmt.Println(s1)

	s2 := append(s, 4)
	fmt.Println(s)
	fmt.Println(s1)
	fmt.Println(s2)
}
