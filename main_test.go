package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	log "github.com/cihub/seelog"
	"github.com/robfig/cron"
)

// Login 模拟登录
func TestAny(t *testing.T) {
	logger, err := log.LoggerFromConfigAsFile("conf/seelog.xml")
	if err != nil {
		fmt.Println(err)
		panic("Open seelog.xml failed!")
	}
	fmt.Println("current time: ", time.Now().Format("2006-01-02 15:04:05"))
	c := cron.New()
	//spec := "5 */1 * * * ?"
	spec := "@every 3s"
	c.AddFunc(spec, func() {
		fmt.Println("cron running: ", time.Now().Format("2006-01-02 15:04:05"))
		logger.Info("cron running: ", time.Now().Format("2006-01-02 15:04:05"))
	})
	c.Start()
	// now := time.Now()
	// fmt.Println(now)
	// end := now.Add(time.Second * 20)
	// fmt.Println(end)
	// fmt.Println(end.Sub(now))
	// timer := time.NewTimer(end.Sub(now))
	// <-timer.C
	fmt.Printf("Server(Listen on %d) is starting...\n", 8080)
	http.ListenAndServe(":8080", nil) //设置监听的端口
}
