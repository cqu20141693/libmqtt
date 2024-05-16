package scheduler

import (
	"log"
	"testing"
	"time"
)

func TestTimeAfter(t *testing.T) {

	log.Println(time.Now())
	<-time.After(5 * time.Second)
	log.Println(time.Now())

}

func TestTimerSelect(t *testing.T) {

	timer := time.NewTimer(time.Second * 3)
	select {
	case <-timer.C:
		log.Println("timed out")
	}

}

func TestAfterFunc(t *testing.T) {
	log.Println("start", time.Now())
	time.AfterFunc(1*time.Second, func() {
		log.Println("time.After 1S ->", time.Now())
	})

	time.Sleep(2 * time.Second) //等待协程退出
}

func TestTicker(t *testing.T) {
	// 建议用gocron 替代Ticker
	// Ticker 需要自己goruntine 调度
	// gocron 会内部自己调度，并且api 方便

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("ticker...")
	}
}
