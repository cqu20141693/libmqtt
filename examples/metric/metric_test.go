package metric

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"sync"
	"testing"
	"time"
)

func TestGoMetric(t *testing.T) {
	counter := metrics.NewCounter()
	name := "pub"
	err := metrics.Register(name, counter)
	if err != nil {
		_ = fmt.Errorf("metrics.Register %s failed", name)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(counter metrics.Counter) {

		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 10)
			counter.Inc(1)
		}
		wg.Done()
	}(counter)
	wg.Add(1)
	go func(counter metrics.Counter) {

		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 10)
			counter.Inc(1)
		}
		wg.Done()
	}(counter)
	wg.Wait()
	fmt.Printf("counter=%d", counter.Count())
	fmt.Printf("pub=%v", metrics.Get(name))

}
