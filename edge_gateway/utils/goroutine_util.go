package utils

import (
	"fmt"
	"log"
	"runtime"
)

// GoWithRecover 可恢复携程
// 避免因为野生 goroutine panic 导致主进程退出
//
//	@param f
func GoWithRecover(f func()) {
	go func() {
		var err error
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err = fmt.Errorf("panic recovered: %s\n%s", r, buf)
			}
			if err != nil {
				log.Printf("panic: %+v", err)
			}
		}()

		f()
	}()
}
