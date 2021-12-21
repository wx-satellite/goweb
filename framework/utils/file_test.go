package utils

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	fmt.Println(Exists("/home/ha"))
}

func TestIsHiddenDirectory(t *testing.T) {
	fmt.Println(filepath.Base("/home/.name"))
}
