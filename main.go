package main

import (
	"Candy/util"
	"fmt"
	"net/http"
)

func main() {
	defer util.Pause()
	authArray, port := util.InitConfig()
	if len(authArray) > 0 {
		for _, value := range authArray {
			if value.Valid {
				go util.DoTask(value.Phone, value.Passwd)
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
