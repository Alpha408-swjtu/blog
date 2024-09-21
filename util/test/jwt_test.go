package test

import (
	"blog/util"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBase64(t *testing.T) {
	text := "你好，大乔乔"
	cipher := base64.StdEncoding.EncodeToString([]byte(text)) //base64是一种公开透明的转码方式，不是加密算法
	fmt.Println(cipher)
	bs, _ := base64.StdEncoding.DecodeString(cipher)
	if string(bs) != text {
		t.Fail()
	}
}

func TestJWT(t *testing.T) {
	secret := "123456"
	header := util.DefautHeader
	payload := util.JwtPayload{
		ID:          "rj4t49tu49",
		Issue:       "微信",
		Audience:    "王者荣耀",
		Subject:     "购买道具",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(2 * time.Hour).Unix(),
		UserDefined: map[string]any{"name": strings.Repeat("大乔乔", 100)}, //信息量很大时，jwt长度可能会超过4K
	}

	if token, err := util.GenJWT(header, payload, secret); err != nil {
		fmt.Printf("生成json web token失败: %v", err)
	} else {
		fmt.Println(token)
		if _, p, err := util.VerifyJwt(token, secret); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("JWT验证通过。欢迎 %s !\n", p.UserDefined["name"])
		}
	}
}

// go test -v .\util\test\ -run=^TestBase64$ -count=1
// go test -v .\util\test\ -run=^TestJWT$ -count=1
