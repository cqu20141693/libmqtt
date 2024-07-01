package collection

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"sync/atomic"
	"testing"
	"time"
)

func TestApi(t *testing.T) {

	queue := gqueue.New(1)
	queue.Pop()
	fmt.Println("collection empty ,pop return right now")
}

func TestPop(t *testing.T) {
	q := gqueue.New()
	ctx := context.Background()

	// 数据生产者，每隔1秒往队列写数据
	gtimer.SetInterval(ctx, time.Second, func(ctx2 context.Context) {

		v := gtime.Now().String()
		q.Push(v)
		fmt.Println("Push:", v)
	})

	// 3秒后关闭队列
	gtimer.SetTimeout(ctx, 3*time.Second, func(ctx2 context.Context) {
		q.Close()
	})

	// 消费者，不停读取队列数据并输出到终端
	for {
		if v := q.Pop(); v != nil {
			fmt.Println(" Pop:", v)
		} else {
			break
		}
	}
}

func TestList(t *testing.T) {
	list := glist.New(true)
	go func() {
		for {
			time.Sleep(time.Millisecond)
			e := list.PopFront()
			if e != nil {
				fmt.Println(fmt.Sprintf("list pop=%v", e))
			}
		}
	}()
	var index uint32

	for i := 0; i < 10; i++ {
		atomic.AddUint32(&index, 1)
		time.Sleep(time.Millisecond)
		list.PushBack(i)
	}

}
func TestArray(t *testing.T) {
	array := garray.New(true)

	go func() {
		for {
			if value, found := array.PopLeft(); found {
				fmt.Println(fmt.Sprintf("array pop=%v", value))
			}

		}
	}()
	for i := 0; i < 100; i++ {
		array.PushRight(i)
	}

	time.Sleep(time.Second)
}
