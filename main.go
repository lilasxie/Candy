package main

import (
	"Candy/util"
	"fmt"
	"net/http"
	_ "net/http/pprof" //监控性能及GC调优

	log "github.com/cihub/seelog"
)

// tryCatch 实现 try catch.
// 空接口：具有0个方法的接口称为空接口。它表示为interface {}。由于空接口有0个方法，所有类型都实现了空接口
func tryCatch(try func(phone, passwd string), handler func(interface{}), phone string, passwd string) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	try(phone, passwd)
}

func main() {
	defer func() {
		log.Flush()
		if err := recover(); err != nil {
			fmt.Printf("%v", err)
			util.Pause()
		}
	}()
	authArray, port := util.InitConfig()
	if len(authArray) > 0 {
		for _, value := range authArray {
			if value.Valid {
				// go util.DoTask(value.Phone, value.Passwd)
				go tryCatch(util.DoTask, func(err interface{}) {
					fmt.Printf("phone:%s, Open treasure box failed, will be removed from account storage. %v\n", value.Phone, err)
				}, value.Phone, value.Passwd)
			}
		}
	}

	fmt.Printf("Server(Listen on %s) is starting...\n", port)
	port = ":" + port
	err := http.ListenAndServe(port, nil) //设置监听的端口
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
