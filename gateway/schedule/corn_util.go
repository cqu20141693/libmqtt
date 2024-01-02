package schedule

import (
	"github.com/go-co-op/gocron"
	"github.com/robfig/cron"
	"time"
)

var Cron = cron.New()
var GoCron = gocron.NewScheduler(time.Local) // 使用系统的本地时区
