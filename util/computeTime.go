package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// ComputeTime 获取上次打开宝箱时间以及宝箱打开间隔参数
func ComputeTime(accessToken string) (string, int64) {

	url := "https://ibo.candy.one/api/period-task/box"
	payload := strings.NewReader("task_id=1")
	req, _ := http.NewRequest("GET", url, payload)

	req.Header.Add("x-access-token", accessToken)
	req.Header.Add("Cache-Control", "no-cache")

	res, _ := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	respStr := string(body)

	codeObj := gjson.Get(respStr, "code")
	if codeObj.Exists() {
		code := codeObj.Int()
		// code为0说明请求失败了
		if code == 0 {
			errorCodeObj := gjson.Get(respStr, "err.err_code")
			if errorCodeObj.Exists() {
				errorCode := errorCodeObj.Int()
				errorMsg := gjson.Get(respStr, "err.message").String()
				fmt.Printf("ComputeTime failed, code: %d, errorCode: %d, message: %s\n", code, errorCode, errorMsg)
				return "error", errorCode
			}
		} else if code == 1 {
			dataObj := gjson.Get(respStr, "data")
			if dataObj.Exists() {
				taskPeriod := dataObj.Get("task_period").Int()
				lastestRewardedTime := dataObj.Get("lastest_rewarded_time").String()
				return lastestRewardedTime, taskPeriod
			}
		} else {
			fmt.Println("ComputeTime, Unknown code: ", code)
			return "error", 0
		}
	}
	return "error", 0
}
