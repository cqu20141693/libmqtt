package schedule

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"log"
	"testing"
	"time"
)

// TestTickerPolicy 测试定时策略
//
//	@param t
func TestTickerPolicy(t *testing.T) {

}

// TestTickerUtilPolicy 测试定时到一个时间点
//
//	@param t
func TestTickerUtilPolicy(t *testing.T) {

	count := 100
	frequency := 1 * time.Second

	duration := time.Duration(count) * frequency

	timer := time.NewTimer(duration)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Println("NewScheduler failed=%v", err)
		// handle error
	}
	counter := 0
	job, err := scheduler.NewJob(
		gocron.DurationJob(
			1*time.Second,
		),
		gocron.NewTask(
			func() {
				counter++
				log.Println("task ", counter)
			}),
	)
	if err != nil {
		return
	}

	scheduler.Start()
	select {
	case <-timer.C:
		err = scheduler.RemoveJob(job.ID())
		if err != nil {
			log.Println("remove job failed=", "jobName")
		}
	}
	err = scheduler.Shutdown()
	if err != nil {
		cclog.SugarLogger.Error("Shutdown schedule failed=", err)
	}

}
