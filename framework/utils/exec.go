package utils

import (
	"os"
	"syscall"
)

// GetExecDirectory 获取当前执行程序的目录
func GetExecDirectory() string {
	file, err := os.Getwd()
	if err != nil {
		return ""
	}
	return file + "/"
}

// CheckProcessExist 当返回true时则表示进程存在
func CheckProcessExist(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))

	if err != nil {
		return false
	}

	return true
}
