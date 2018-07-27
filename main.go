package main

import (
	"Candy/util"
	"fmt"
	"net/http"
)

func main() {
	authArray, port := util.InitConfig()
	if len(authArray) > 0 {
		for _, value := range authArray {
			if value.Valid {
				go util.DoTask(value.Phone, value.Passwd)
			}
		}
	}

	port = ":" + port
	err := http.ListenAndServe(port, nil) //设置监听的端口
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		util.Pause()
	}
	fmt.Println("Start server success, listen on port", port)
}
