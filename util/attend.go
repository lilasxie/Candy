package util

import (
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

// Attend 模拟签到
func Attend(accessToken string) (int64, int64) {
	url := "https://ibo.candy.one/api/ibo-task/clockin-daily"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
				return code, errorCode
			}
		} else if code == 1 {
			return code, 1
		} else {
			return code, 0
		}
	}
	return 0, 0
}
