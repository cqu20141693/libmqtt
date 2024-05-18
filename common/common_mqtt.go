package common

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/go-co-op/gocron/v2"
	"github.com/goiiot/libmqtt"
	"github.com/goiiot/libmqtt/domain"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	serve "github.com/goiiot/libmqtt/edge_gateway/initialize/server"
	"github.com/goiiot/libmqtt/edge_gateway/mqtt"
	"github.com/rcrowley/go-metrics"
	"sync"
	"sync/atomic"
	"time"
)

var once = sync.Once{}
var connOk = "connectSuccess"
var connFail = "connectFailed"
var sum = "sum"
var connOkCounter = metrics.NewCounter()
var connFailCounter = metrics.NewCounter()
var sumCounter = metrics.NewCounter()

func initMetric() {
	once.Do(func() {
		err := metrics.Register(connOk, connOkCounter)
		if err != nil {
			cclog.SugarLogger.Error("initMetric connectSuccess failed")
		}

		err = metrics.Register(sum, connFailCounter)
		if err != nil {
			cclog.SugarLogger.Error("initMetric sum failed")
		}
		err = metrics.Register(connFail, sumCounter)
		if err != nil {
			cclog.SugarLogger.Error("initMetric connectFailed failed")
		}
	})

}

// MqttConnect 并发连接MQTT server
//
//	@param count 并发数量
//	@param deviceFormat 设备格式：eg "direct%d"
//	@param protoVersion
func MqttConnect(count int, connInfoFunc ConnectInfoFunc) {
	infos := make([]*domain.MqttClientAddInfo, 0, count)
	ret := struct {
		success atomic.Int32
		failed  atomic.Int32
	}{}
	for i := 0; i < count; i++ {
		server, clientId, username, password, keepalive, protoVersion, duration := connInfoFunc(i)
		info := domain.NewMqttClientAddInfoWithVersion(server, clientId, username, password, keepalive, byte(protoVersion))
		infos = append(infos, info)
		go func(addInfo *domain.MqttClientAddInfo) {

			client, err := mqtt.CreatClient(addInfo)
			if err != nil {
				_ = fmt.Errorf("creat g platform mqtt client failed, %v", info)
				ret.failed.Add(1)

				return
			}
			ret.success.Add(1)
			if duration > 0 {
				time.Sleep(duration)
			}

			client.Wait()
		}(info)

	}

	fmt.Println(fmt.Sprintf("success=%d,failed=%d", ret.success.Load(), ret.failed.Load()))
}

//  @param  index 索引
//  @return string server
//  @return string clientId
//  @return string username
//  @return string password
//  @return int64 keepalive

type ConnectInfoFunc func(index int) (string, string, string, string, int64, libmqtt.ProtoVersion, time.Duration)

// eg: publishInfo

// @param clientId
// @return DirectTelemetryFunc 遥测数据生成器
// @return string topic
// @return int fasong  ： -1 永久； 0 一次； n次
// @return libmqtt.QosLevel qos
// @return time.Duration 发送频率
type PublishInfoFunc func(clientId string) (TelemetryFunc, string, int, libmqtt.QosLevel, time.Duration)

type TelemetryFunc func(deviceId string) []map[string]interface{}

// MqttPublish 连接mqtt并publish 消息
//
//	@param count 连接的mqtt 数量
//	@param connInfoFunc 连接信息
//	@param pubInfoFunc 推送信息
func MqttPublish(count int, connInfoFunc ConnectInfoFunc, pubInfoFunc PublishInfoFunc) {

	infos := make([]*domain.MqttClientAddInfo, 0, count)
	group := sync.WaitGroup{}
	scheduler, _ := gocron.NewScheduler(gocron.WithLocation(time.Local),
		gocron.WithGlobalJobOptions(gocron.WithSingletonMode(gocron.LimitModeReschedule)))
	group.Add(count)
	for i := 0; i < count; i++ {

		server, clientId, username, password, keepalive, protoVersion, duration := connInfoFunc(i)
		telemetryFunc, topic, pubCount, qosLevel, frequency := pubInfoFunc(clientId)
		initMetric()
		info := domain.NewMqttClientAddInfoWithVersion(server, clientId, username, password, keepalive, byte(protoVersion))
		infos = append(infos, info)
		go func(addInfo *domain.MqttClientAddInfo) {

			client, err := mqtt.CreatClient(addInfo)
			if err != nil {
				cclog.SugarLogger.Error("creat g platform mqtt client failed, %v", info)
				connFailCounter.Inc(1)
				group.Done()
				return
			}
			info.SetClient(client)
			connOkCounter.Inc(1)
			// 等待连接完成

			time.Sleep(time.Second * 2)
			publishBySleep(telemetryFunc, info, client, topic, qosLevel, pubCount, frequency)
			// 存在发送任务为调度的情况，需要了解框架，保证调度任务都被执行
			//publishByScheduler(telemetryFunc, info, client, topic, qosLevel, time.Duration(pubCount)*frequency, scheduler, frequency, clientId)
			if duration > 0 {
				time.Sleep(duration)
			} else {
				time.Sleep(time.Second * 3)
			}

			if client.PubMetric != nil {

				count := client.PubMetric.Count()
				cclog.SugarLogger.Info(fmt.Sprintf("device: %s, publish: %d", info.ClientID, count))
				sumCounter.Inc(count)
			}
			group.Done()
		}(info)

	}

	printTimer(scheduler, infos)
	scheduler.Start()
	group.Wait()
	_ = scheduler.Shutdown()
	cclog.SugarLogger.Info(fmt.Sprintf("success=%d,failed=%d,sum=%d", connOkCounter.Count(), connFailCounter.Count(), sumCounter.Count()))
}

func printTimer(scheduler gocron.Scheduler, infos []*domain.MqttClientAddInfo) {
	_, err := scheduler.NewJob(gocron.DurationJob(time.Second*30), gocron.NewTask(func() {
		for _, info := range infos {
			info.PrintClientMetric()
		}
	}))
	if err != nil {
		cclog.SugarLogger.Error(err)
	}
}

func publishBySleep(telemetryFunc TelemetryFunc, info *domain.MqttClientAddInfo, client libmqtt.Client, topic string, qosLevel libmqtt.QosLevel, pubCount int, frequency time.Duration) {
	publishFunc := func() {
		telemetryPkt := telemetryFunc(info.ClientID)
		for _, pkt := range telemetryPkt {
			telemetryPktBytes, err := sonic.Marshal(pkt)
			if err != nil {
				cclog.SugarLogger.Errorf("sonic marchal faild: %v", err)
				return
			}
			client.Publish(&libmqtt.PublishPacket{
				TopicName: topic,
				Qos:       qosLevel,
				Payload:   telemetryPktBytes,
			})
		}

	}
	if pubCount < 0 {
		for true {

			time.Sleep(frequency)
			publishFunc()
			select {
			case <-serve.MainCtx.Done():
				cclog.Warn("job complete receive SIGTERM ....")
				return
			default:
				continue
			}
		}
	} else if pubCount == 0 {
		publishFunc()
	} else {

		for i := 0; i < pubCount; i++ {
			time.Sleep(frequency)
			publishFunc()
			select {
			case <-serve.MainCtx.Done():
				cclog.Warn("job complete receive SIGTERM ....")
				return
			default:
				continue
			}
		}
	}
}

// publishByScheduler 通过调度器发送消息
// 需要优化解决任务调度调度完成，不能出现部分任务未调度情况
//
//	@param telemetryFunc
//	@param info
//	@param client
//	@param topic
//	@param qosLevel
//	@param duration
//	@param scheduler
//	@param frequency
//	@param clientId
func publishByScheduler(telemetryFunc TelemetryFunc, info *domain.MqttClientAddInfo, client libmqtt.Client, topic string, qosLevel libmqtt.QosLevel, duration time.Duration, scheduler gocron.Scheduler, frequency time.Duration, clientId string) {
	publishFunc := func() {
		telemetryPkt := telemetryFunc(info.ClientID)
		telemetryPktBytes, err := sonic.Marshal(telemetryPkt)
		if err != nil {
			cclog.SugarLogger.Errorf("sonic marchal faild: %v", err)
		}
		client.Publish(&libmqtt.PublishPacket{
			TopicName: topic,
			Qos:       qosLevel,
			Payload:   telemetryPktBytes,
		})
	}
	// 支持策略定时发送，定时发送一段时间，发送一次
	if duration > 0 {

		job, err := scheduler.NewJob(gocron.DurationJob(frequency),
			gocron.NewTask(publishFunc),
			//gocron.WithStartAt(gocron.WithStartImmediately()),
		)
		if err != nil {
			cclog.SugarLogger.Errorf("schedule job failed=%s %v", clientId, err)
		}

		select {
		case <-time.After(duration):
			cclog.SugarLogger.Warnf("job complete duration=%d arrived.", duration)
			err := scheduler.RemoveJob(job.ID())
			if err != nil {
				cclog.SugarLogger.Errorf("schedule remove job failed=%s %v", clientId, err)
			}
		case <-serve.MainCtx.Done():
			cclog.SugarLogger.Warn("job complete receive SIGTERM ....")
			err := scheduler.RemoveJob(job.ID())
			if err != nil {
				cclog.SugarLogger.Errorf("schedule remove job failed=%s %v", clientId, err)
			}
		}

	} else if duration == 0 {
		t := time.Now().Add(frequency)
		_, _ = scheduler.NewJob(gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(t)),
			gocron.NewTask(publishFunc),
		)

	} else {
		job, err := scheduler.NewJob(gocron.DurationJob(frequency),
			gocron.NewTask(publishFunc),
			//gocron.WithStartAt(gocron.WithStartImmediately()),
		)
		if err != nil {
			cclog.SugarLogger.Errorf("schedule job failed=%s %v", clientId, err)
		}
		select {
		case <-serve.MainCtx.Done():
			cclog.Warn("job complete receive SIGTERM ....")
			err := scheduler.RemoveJob(job.ID())
			if err != nil {
				cclog.SugarLogger.Errorf("schedule remove job failed=%s %v", clientId, err)
			}
		}
	}
}
