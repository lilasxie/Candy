package util

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/cihub/seelog"
	"github.com/robfig/cron"
)

var logger = getLogger()

// getLogger 获取日志句柄
func getLogger() log.LoggerInterface {
	//加载配置文件
	logger, err := log.LoggerFromConfigAsFile("conf/seelog.xml")
	if err != nil {
		panic("Open seelog.xml failed!")
	}
	return logger
}

// keepGroutineAlive 设置心跳，防止进程假死
func keepGroutineAlive(phone string) *cron.Cron {
	taskCron := cron.New()
	//spec := "0 */5 * * * ?"
	spec := "@every 30m"
	taskCron.AddFunc(spec, func() {
		logger.Infof("Keep goroutine alive for phone : %s, current time: %s\n", phone, time.Now().Format("2006-01-02 15:04:05"))
	})
	taskCron.Start()
	return taskCron
}

// sleepAfterSuccess 休眠直到可以再次打开宝箱
func sleepAfterSuccess(nextRewardedTimeUnix, curUnixTimeUnix int64, phone string) {
	// 生成一个1-20分钟的随机加时
	rand.Seed(time.Now().UnixNano())
	extTimeUnix := rand.Int63n(1199) + 1
	nextOpenTimeUnix := nextRewardedTimeUnix + extTimeUnix
	//需等待时长
	sleepTimeUnix := time.Duration(nextRewardedTimeUnix-curUnixTimeUnix+extTimeUnix) * time.Second
	logger.Infof("phone:%s, Treasure box will be opened on %v\n", phone, time.Unix(nextOpenTimeUnix, 0).In(time.Local))
	fmt.Printf("phone:%s, Treasure box will be opened on %v\n", phone, time.Unix(nextOpenTimeUnix, 0).In(time.Local))
	// time.Sleep(sleepTimeUnix)
	timer := time.NewTimer(sleepTimeUnix)
	<-timer.C
}

//DoTask 模拟执行各个任务
func DoTask(phone, passwd string) {
	pingCron := keepGroutineAlive(phone)
	defer pingCron.Stop()
	accessToken := Login(phone, passwd)
	// 签到任务下线，关闭签到功能(pS:后台接口并未关闭, 但是大佬说还是要关闭签到功能)
	// 假设每次启动都还未签到，默认上一次签到为昨天
	// lastAttendanceDay := time.Now().AddDate(0, 0, -1).Day()
	for {
		if accessToken != "" {
			logger.Infof("phone:%s, accessToken = %s\n", phone, accessToken)
			// 签到任务
			/*currentDay := time.Now().Day()
			logger.Infof("phone:%s, currentDay:%d, lastAttendanceDay:%d\n", phone, currentDay, lastAttendanceDay)
			if currentDay != lastAttendanceDay {
				resultCode, resultErrorCode := Attend(accessToken)
				if resultCode == 1 {
					// 更新最近一次签到时间
					lastAttendanceDay = currentDay
					logger.Infof("phone:%s, Attend success\n", phone)
				} else if resultCode == 0 && resultErrorCode == 10008 {
					// 说明accessToken过期，需要重新登录
					logger.Infof("phone:%s, Attend failed for accessToken invalid, request new accessToken......\n", phone)
					accessToken = Login(phone, passwd)
				} else if resultCode == 0 && resultErrorCode == 400021 {
					lastAttendanceDay = currentDay
					logger.Infof("phone:%s, Attend failed for today has attended\n", phone)
				} else {
					logger.Infof("phone:%s, Attend failed, code:%d, err_code:%d\n", phone, resultCode, resultErrorCode)
				}
			}*/

			// 开宝箱任务
			lastestRewardedTimeStr, taskPeriod := ComputeTime(accessToken)
			if lastestRewardedTimeStr != "error" && taskPeriod > 0 {
				// logger.Infof("phone:%s, ComputeTime success\n", phone)
				lastestRewardedTime, err := time.Parse(time.RFC3339, lastestRewardedTimeStr)
				if err == nil {
					logger.Infof("phone:%s, Last rewarded time: %v\n", phone, lastestRewardedTime.In(time.Local))
					// fmt.Printf("phone:%s, Last rewarded time: %v\n", phone, lastestRewardedTime.In(time.Local))
					lastestRewardedTimeUnix := lastestRewardedTime.Unix()
					nextRewardedTimeUnix := lastestRewardedTimeUnix + 60*taskPeriod
					logger.Infof("phone:%s, Next rewarded time: %v\n", phone, time.Unix(nextRewardedTimeUnix, 0).In(time.Local))
					curUnixTimeUnix := time.Now().Unix()
					logger.Infof("phone:%s, Current time: %v\n", phone, time.Now().In(time.Local))
					if curUnixTimeUnix > nextRewardedTimeUnix {
						// 说明又可以开宝箱了
						resultCode, resultErrorCode := OpenTreasureAction(accessToken)
						if resultCode == 1 {
							logger.Infof("phone:%s, Open treasure box success, waiting for next\n", phone)
							fmt.Printf("phone:%s, Open treasure box success on %s, waiting for next\n", phone, time.Now().Format("2006-01-02 15:04:05"))
							curUnixTimeUnix = time.Now().Unix()
							nextRewardedTimeUnix = curUnixTimeUnix + 60*taskPeriod
							sleepAfterSuccess(nextRewardedTimeUnix, curUnixTimeUnix, phone)
						} else if resultCode == 0 && resultErrorCode == 10008 {
							// 说明accessToken过期，需要重新登录
							logger.Infof("phone:%s, Open treasure box failed for accessToken invalid, request new accessToken......\n", phone)
							accessToken = Login(phone, passwd)
						} else {
							logger.Infof("phone:%s, Open treasure box failed, code:%d, err_code:%d\n", phone, resultCode, resultErrorCode)
						}
					} else {
						logger.Infof("phone:%s, Open treasure box failed for time has not reached\n", phone)
						sleepAfterSuccess(nextRewardedTimeUnix, curUnixTimeUnix, phone)
					}
				} else {
					logger.Info(err)
				}
			} else {
				logger.Infof("phone:%s, accessToken invalid, request new accessToken......\n", phone)
				accessToken = Login(phone, passwd)
			}
		} else {
			logger.Infof("phone:%s, Invalid accessToken\n", phone)
			logger.Infof("phone:%s, will be removed from account storage\n", phone)
			// fmt.Printf("phone:%s, will be removed from account storage\n", phone)
			break
		}
	}
}
