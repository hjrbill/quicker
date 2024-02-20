package util

import (
	"path"
	"runtime"
)

var (
	RootPath string //项目根目录
)

func init() {
	RootPath = path.Dir(GetCurrentPath()) + "/../" //项目根目录
}

// GetCurrentPath 获取调用方所在 go 代码的路径
func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
