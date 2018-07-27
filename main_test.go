package main

import (
	"log"
	"testing"

	"github.com/robfig/cron"
)

// Login 模拟登录
func TestAny(t *testing.T) {
	i := 0
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	c.Start()

	select {}
}
