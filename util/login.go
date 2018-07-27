package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// Login 模拟登录
func Login(phone, passwd string) string {
	url := "https://ibo.candy.one/api/passport/password-login"
	payload := strings.NewReader("country_code=cn&password=" + passwd + "&phone=%2B86" + phone)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	res, _ := http.DefaultClient.Do(req)
	// TODO 验证登录是否成功
	if res != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	respStr := string(body)

	codeObj := gjson.Get(respStr, "code")
	if codeObj.Exists() {
		code := codeObj.Int()
		if code == 1 {
			accessToken := gjson.Get(respStr, "data.access_token")
			if accessToken.Exists() {
				return accessToken.String()
			}
			fmt.Printf("phone:%s, 解析accessToken异常\n", phone)
		} else {
			msg := gjson.Get(respStr, "err.message").String()
			fmt.Printf("phone:%s, 登录失败,%s\n", phone, msg)
		}
	} else {
		fmt.Printf("phone:%s, 解析accessToken异常\n", phone)
	}
	return ""
}
