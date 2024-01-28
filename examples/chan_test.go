package examples

import (
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"sync"
	"testing"
	"time"
)

// chan 缓冲的长度和容量两个小知识点
// 同时测试chan 死锁
func TestChanDeadLock(t *testing.T) {
	defer func() {
		a := recover()
		fmt.Println(a)
	}()
	ch := make(chan string, 3)
	ch <- "jiyik.com"
	ch <- "onmpw.com"
	fmt.Println("容量是 ", cap(ch))
	fmt.Println("长度是 ", len(ch))
	fmt.Println("读取值 ", <-ch)
	fmt.Println("新的长度是 ", len(ch))
	ch <- "test1"
	fmt.Println("新的长度是 ", len(ch))
	ch <- "test2"
	fmt.Println("新的长度是 ", len(ch))
	// channel 阻塞，出现死锁
	ch <- "test3"
}

func TestSelectChan(t *testing.T) {

}

// TestChanProducerAndConsumer
// 测试chan生产者消费者模式
//  @param t
func TestChanProducerAndConsumer(t *testing.T) {
	eventCh := make(chan string)
	service1 := Service1{
		event: eventCh,
	}
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		for {
			time.Sleep(time.Second * 3)
			service1.send(time.Now().String())
		}
		group.Done()
	}()
	service2 := Service2{
		event: eventCh,
	}
	group.Add(1)
	go func() {
		service2.receive()
		group.Done()
		service2.receive2()
		group.Done()
	}()
	group.Wait()
}

type Service1 struct {
	event chan string
}

func (s *Service1) send(data string) {
	s.event <- data
}

type Service2 struct {
	event chan string
}

func (s *Service2) receive() {
	select {
	case e, more := <-s.event:
		{
			if !more {
				return
			}
			cclog.Info("receive data=", e)
		}

	}
}

func (s *Service2) receive2() {
	for true {
		select {
		case e, more := <-s.event:
			{
				if !more {
					return
				}
				cclog.Info("receive data=", e)
			}

		}
	}

}
