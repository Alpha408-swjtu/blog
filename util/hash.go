package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(text string) string {
	md5 := md5.New()
	md5.Write([]byte(text))
	digest := md5.Sum(nil)            //md5哈希的结果是128bit
	return hex.EncodeToString(digest) //十六进制编码之后是128/4=32个字符
}
