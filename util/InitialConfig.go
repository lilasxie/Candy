package util

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tidwall/gjson"
)

// Auth 账号密码
type Auth struct {
	Phone  string
	Passwd string
	Valid  bool
}

// InitConfig 初始化配置
func InitConfig() ([]Auth, string) {
	file, err := os.Open("conf/conf.json")
	defer file.Close()
	var authArray []Auth
	if err != nil {
		fmt.Println("Initialize config failed, cause:", err)
		return authArray, "8888"
	}
	conf, _ := ioutil.ReadAll(file)
	confStr := string(conf)
	authObj := gjson.Get(confStr, "auth")
	if authObj.Exists() {
		for _, value := range authObj.Array() {
			phone := value.Get("phone").String()
			passwd := value.Get("passwd").String()
			valid := value.Get("valid").Bool()
			authArray = append(authArray, Auth{phone, passwd, valid})
		}
	}

	port := "8888"
	portObj := gjson.Get(confStr, "server.port")
	if portObj.Exists() {
		port = portObj.String()
	}
	fmt.Println("Initialize config success")
	return authArray, port
}
