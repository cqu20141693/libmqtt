package file

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var osType string
var path string

const WINDOWS = "windows"

func init() {
	osType = runtime.GOOS
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
}

func MkdirIfNecessary(createDir string) (err error) {
	return os.MkdirAll(createDir, os.ModePerm)
}

func GetCurrentPath() string {
	dir, err := os.Getwd() //当前的目录
	if err != nil {
		dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Println("can not get current path")
		}
	}
	return dir
}
