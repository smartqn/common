package util

import (
	"path"
	"runtime"
	"strings"
	"sync"
)

var _level = -1
var once sync.Once

func curFile(addLevel int) string {
	once.Do(func() {
		var filename string
		for i := 0; i < 20; i++ {
			_, filename, _, _ = runtime.Caller(i)
			if strings.HasSuffix(filename, "config/config.go") {
				_level = i + 1
				break
			}
		}

	})
	_, filename, _, _ := runtime.Caller(_level + addLevel)
	return filename
}

// 获取调用者的当前文件名
func CurFile() string {
	return curFile(1)
}

// 获取调用者的当前文件DIR
func CurDir() string {
	return path.Dir(curFile(1))
}
