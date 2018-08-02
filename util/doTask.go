package util

import (
	"fmt"
	"math/rand"
	"time"
)

//DoTask 模拟执行各个任务
func DoTask(phone, passwd string) {
	accessToken := Login(phone, passwd)
	// 签到任务下线，关闭签到功能(pS:后台接口并未关闭, 但是大佬说还是要关闭签到功能)
	// 假设每次启动都还未签到，默认上一次签到为昨天
	// lastAttendanceDay := time.Now().AddDate(0, 0, -1).Day()
	for {
		fmt.Println("**************************************" + phone + "*****************************************")
		if accessToken != "" {
			fmt.Printf("phone:%s, accessToken = %s\n", phone, accessToken)
			// 签到任务
			/*currentDay := time.Now().Day()
			fmt.Printf("phone:%s, currentDay:%d, lastAttendanceDay:%d\n", phone, currentDay, lastAttendanceDay)
			if currentDay != lastAttendanceDay {
				resultCode, resultErrorCode := Attend(accessToken)
				if resultCode == 1 {
					// 更新最近一次签到时间
					lastAttendanceDay = currentDay
					fmt.Printf("phone:%s, Attend success\n", phone)
				} else if resultCode == 0 && resultErrorCode == 10008 {
					// 说明accessToken过期，需要重新登录
					fmt.Printf("phone:%s, Attend failed for accessToken invalid, request new accessToken......\n", phone)
					accessToken = Login(phone, passwd)
				} else if resultCode == 0 && resultErrorCode == 400021 {
					lastAttendanceDay = currentDay
					fmt.Printf("phone:%s, Attend failed for today has attended\n", phone)
				} else {
					fmt.Printf("phone:%s, Attend failed, code:%d, err_code:%d\n", phone, resultCode, resultErrorCode)
				}
			}*/

			// 开宝箱任务
			lastestRewardedTimeStr, taskPeriod := ComputeTime(accessToken)
			if lastestRewardedTimeStr != "error" && taskPeriod > 0 {
				// fmt.Printf("phone:%s, ComputeTime success\n", phone)
				lastestRewardedTime, err := time.Parse(time.RFC3339, lastestRewardedTimeStr)
				if err == nil {
					fmt.Printf("phone:%s, Last rewarded time: %v\n", phone, lastestRewardedTime.In(time.Local))
					lastestRewardedTimeUnix := lastestRewardedTime.Unix()
					nextRewardedTimeUnix := lastestRewardedTimeUnix + 60*taskPeriod
					fmt.Printf("phone:%s, Next rewarded time: %v\n", phone, time.Unix(nextRewardedTimeUnix, 0).In(time.Local))
					curUnixTimeUnix := time.Now().Unix()
					fmt.Printf("phone:%s, Current time: %v\n", phone, time.Now().In(time.Local))
					if curUnixTimeUnix > nextRewardedTimeUnix {
						// 说明又可以开宝箱了
						resultCode, resultErrorCode := OpenTreasureAction(accessToken)
						if resultCode == 1 {
							fmt.Printf("phone:%s, Open treasure box success, waiting for next\n", phone)
						} else if resultCode == 0 && resultErrorCode == 10008 {
							// 说明accessToken过期，需要重新登录
							fmt.Printf("phone:%s, Open treasure box failed for accessToken invalid, request new accessToken......\n", phone)
							accessToken = Login(phone, passwd)
						} else {
							fmt.Printf("phone:%s, Open treasure box failed, code:%d, err_code:%d\n", phone, resultCode, resultErrorCode)
						}
					} else {
						// 生成一个1-20分钟的随机加时
						rand.Seed(time.Now().UnixNano())
						extTimeUnix := rand.Int63n(1199) + 1
						nextOpenTimeUnix := nextRewardedTimeUnix + extTimeUnix
						//需等待时长
						sleepTimeUnix := time.Duration(nextRewardedTimeUnix-curUnixTimeUnix+extTimeUnix) * time.Second
						fmt.Printf("phone:%s, Open treasure box failed for time has not reached, will be opened on %v\n", phone, time.Unix(nextOpenTimeUnix, 0).In(time.Local))
						time.Sleep(sleepTimeUnix)
					}
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Printf("phone:%s, accessToken invalid, request new accessToken......\n", phone)
				accessToken = Login(phone, passwd)
			}
		} else {
			fmt.Printf("phone:%s, Invalid accessToken\n", phone)
			fmt.Printf("phone:%s, will be removed from account storage\n", phone)
			break
		}
		fmt.Println("**************************************" + phone + "*****************************************")
	}
}
