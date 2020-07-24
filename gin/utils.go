package gin

import (
	"path"
	"reflect"
	"runtime"
)

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}

	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)

	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}

	return finalPath
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

//获取调用函数信息
func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}

type H map[string]interface{}
