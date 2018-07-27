package util

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// OpenTreasureAction 模拟打开宝箱
func OpenTreasureAction(accessToken string) (int64, int64) {

	url := "https://ibo.candy.one/api/period-task/do"
	payload := strings.NewReader("task_id=1")
	req, _ := http.NewRequest("POST", url, payload)

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
		// code为0说明宝箱打开失败
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
