package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

//JWT: Json Web Token

const (
	JWT_SECRET = "123456" //密钥必须在server端保管好，绝对不能外泄
)

var (
	DefautHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

type JwtHeader struct {
	Algo string `json:"alg"` //哈希算法，默认为HMAC SHA256(写为 HS256)
	Type string `json:"typ"` //令牌(token)类型，统一写为JWT
}

type JwtPayload struct {
	ID          string         `json:"jti"` //JWT ID用于标识该JWT
	Issue       string         `json:"iss"` //发行人。比如微信
	Audience    string         `json:"aud"` //受众人。比如王者荣耀
	Subject     string         `json:"sub"` //主题
	IssueAt     int64          `json:"iat"` //发布时间,精确到秒
	NotBefore   int64          `json:"nbf"` //在此之前不可用,精确到秒
	Expiration  int64          `json:"exp"` //到期时间,精确到秒
	UserDefined map[string]any `json:"ud"`  //用户自定义的其他字段
}

func GenJWT(header JwtHeader, payload JwtPayload, secret string) (string, error) {
	var part1, part2, signature string
	//header转成json，然后进行Base64编码
	if bs1, err := json.Marshal(header); err != nil {
		return "", err
	} else {
		part1 = base64.RawURLEncoding.EncodeToString(bs1) //这里没有使用StdEncoding，RawURLEncoding的结果中不会包含=+/等url中的特殊字符
	}
	//payload转成json，然后进行Base64编码
	if bs2, err := json.Marshal(payload); err != nil {
		return "", err
	} else {
		part2 = base64.RawURLEncoding.EncodeToString(bs2)
	}
	//基于sha256的哈希认证算法。任意长度的字符串，经过sha256之后长度都变成了256 bits
	h := hmac.New(sha256.New, []byte(secret))
	//signature = HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload),secret)
	h.Write([]byte(part1 + "." + part2))
	signature = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return part1 + "." + part2 + "." + signature, nil
}

func VerifyJwt(token string, secret string) (*JwtHeader, *JwtPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("token是%d部分", len(parts))
	}
	//进行哈希签名的验证
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(parts[0] + "." + parts[1]))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if signature != parts[2] { //验证失败
		return nil, nil, fmt.Errorf("验证失败")
	}

	var part1, part2 []byte
	var err error
	if part1, err = base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return nil, nil, fmt.Errorf("header Base64反解失败")
	}
	if part2, err = base64.RawURLEncoding.DecodeString(parts[1]); err != nil {
		return nil, nil, fmt.Errorf("payload Base64反解失败")
	}

	var header JwtHeader
	var payload JwtPayload
	if err = json.Unmarshal(part1, &header); err != nil {
		return nil, nil, fmt.Errorf("header json反解失败")
	}
	if err = json.Unmarshal(part2, &payload); err != nil {
		return nil, nil, fmt.Errorf("payload json反解失败")
	}
	return &header, &payload, nil
}
