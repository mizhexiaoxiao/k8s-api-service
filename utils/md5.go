package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

func CreateSign(appKey, appSecret string) string {
	appKey = "value"
	appSecret = "twsxwflergjiweqwsq"
	encryptStr := appKey + "&ts=xxx"

	// 自定义验证规则
	sn := MD5(encryptStr + appSecret)
	return sn
}
