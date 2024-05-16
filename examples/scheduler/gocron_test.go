package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"sync"
	"testing"
	"time"
)

func TestEveryImmediately(t *testing.T) {
	s, _ := gocron.NewScheduler(gocron.WithLocation(time.Local))
	semaphore := make(chan bool)
	fn := func() { semaphore <- true }

	// WithStartImmediately exec
	_, err := s.NewJob(
		gocron.DurationJob(
			2*time.Second,
		),
		gocron.NewTask(fn),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	require.NoError(t, err)
	s.Start()
	select {
	case <-time.After(time.Second * 3):
		err = s.Shutdown()
		if err != nil {
			t.Fatal("Shutdown scheduler failed")
		}
		t.Fatal("job did not run immediately")
	case <-semaphore:
		// test passed
		err = s.Shutdown()
		if err != nil {
			t.Fatal("Shutdown scheduler failed")
		}
	}
}

func TestEveryWaitForSchedule(t *testing.T) {
	s, _ := gocron.NewScheduler(gocron.WithLocation(time.Local))
	semaphore := make(chan bool)
	fn := func() { semaphore <- true }

	//  默认调度时间到执行
	_, err := s.NewJob(
		gocron.DurationJob(
			2*time.Second,
		),
		gocron.NewTask(fn),
	)
	require.NoError(t, err)
	s.Start()
	select {
	case <-time.After(time.Second * 3):
		err = s.Shutdown()
		if err != nil {
			t.Fatal("Shutdown scheduler failed")
		}
		t.Fatal("job did not run immediately")
	case <-semaphore:
		// test passed
		err = s.Shutdown()
		if err != nil {
			t.Fatal("Shutdown scheduler failed")
		}
	}
}

func TestEverySingletonMode(t *testing.T) {
	// 全局配置WithSingletonMode
	s, _ := gocron.NewScheduler(gocron.WithLocation(time.Local),
		gocron.WithGlobalJobOptions(gocron.WithSingletonMode(gocron.LimitModeReschedule)))
	// 如果之前的任务尚未完成，单例模式将阻止新任务启动
	//  默认调度时间到执行
	_, err := s.NewJob(
		gocron.DurationJob(
			2*time.Second,
		),
		gocron.NewTask(task),
		// 单任务配置
		//gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		return
	}
	// 同步启动，阻塞进程
	s.Start()

	select {
	case <-time.NewTimer(time.Second * 3).C:
		_ = s.Shutdown()
	}

}

func TestCron(t *testing.T) {
	// 全局配置WithSingletonMode
	s, _ := gocron.NewScheduler(gocron.WithLocation(time.Local),
		gocron.WithGlobalJobOptions(gocron.WithSingletonMode(gocron.LimitModeReschedule)))

	// 标准的crontab格式，最小单位是分
	_, _ = s.NewJob(
		gocron.CronJob(
			// standard cron tab parsing
			"* * * * * *",
			true,
		),
		gocron.NewTask(task),
	)
	// 最小单位是秒的crontab表达式
	s.Start()
	select {
	case <-time.NewTimer(time.Second * 3).C:
		_ = s.Shutdown()
	}
}

func TestSchedulerOnce(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		// 全局配置WithSingletonMode
		s, _ := gocron.NewScheduler(gocron.WithLocation(time.Local),
			gocron.WithGlobalJobOptions(gocron.WithSingletonMode(gocron.LimitModeReschedule)))
		now := time.Now()
		semaphore := make(chan bool)
		nextMinuteTime := now.Add(time.Minute)
		startAt := fmt.Sprintf("%02d:%02d:%02d", nextMinuteTime.Hour(), nextMinuteTime.Minute(), nextMinuteTime.Second())
		log.Println("nextTime=", startAt)
		dayJob, err := s.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(
					gocron.NewAtTime(uint(nextMinuteTime.Hour()), uint(nextMinuteTime.Minute()), uint(nextMinuteTime.Second())),
				),
			),
			gocron.NewTask(
				func() {
					log.Println("exec")
					semaphore <- true
				},
			),
		)
		require.NoError(t, err)

		s.Start()

		select {
		case <-time.After(time.Minute):
			_ = s.Shutdown()
			nextRun, _ := dayJob.NextRun()
			assert.Equal(t, now.Add(1*time.Minute), nextRun)
		case <-semaphore:
			_ = s.Shutdown()
			t.Error("job ran even though scheduled in future")
		}
		time.Sleep(2 * time.Second)
		wg.Done()
	}()
	wg.Wait()
	log.Println("success")
}

func task() {
	log.Println("hello")
}
