package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const COOKIE_NAME = "auth"

func Login(ctx *gin.Context) {
	//登录成功,给前端返回一个Set-Cookie
	uid := "8"
	payload := JwtPayload{UserDefined: map[string]any{"uid": uid}}
	token, _ := GenJWT(DefautHeader, payload, JWT_SECRET)
	fmt.Println("jwt token", token)

	//response header里会有一条 Set-Cookie: auth=xxx; other_key=other_value，浏览器后续请求会自动把同域名下的cookie再放到request header里来，即request header里会有一条Cookie: auth=xxx; other_key=other_value
	ctx.SetCookie(
		COOKIE_NAME, //cookie的name
		token,       //cookie的value
		86400*7,     //cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
		"/",         //cookie存放目录
		"localhost", //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问
		false,       //是否只能通过https访问
		true,        //是否允许别人通过js获取自己的cookie，设为false防止XSS攻击
	)
	ctx.String(http.StatusOK, "success")
}

func getUidFromCookie1(ctx *gin.Context) string {
	//http协议里没有cookie这个概念，cookie本质上是header里的一对KV
	for _, cookie := range strings.Split(ctx.Request.Header.Get("cookie"), ";") {
		arr := strings.Split(cookie, "=")
		key := strings.TrimSpace(arr[0])
		value := strings.TrimSpace(arr[1])
		if key == COOKIE_NAME {
			if _, payload, err := VerifyJwt(value, JWT_SECRET); err == nil {
				if uid, exists := payload.UserDefined["uid"]; exists {
					return uid.(string)
				}
			}
		}
	}
	return ""
}
func getUidFromCookie2(ctx *gin.Context) string {
	//直接使用封装好的Request.Cookies()
	for _, cookie := range ctx.Request.Cookies() {
		fmt.Println(cookie.Name, cookie.Value)
		if cookie.Name == COOKIE_NAME {
			if _, payload, err := VerifyJwt(cookie.Value, JWT_SECRET); err == nil {
				if uid, exists := payload.UserDefined["uid"]; exists {
					return uid.(string)
				}
			}
		}
	}
	return ""
}

func Account(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "account.html", gin.H{"uid": ctx.GetString("uid")})
}

func Manage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "manage.html", gin.H{"uid": ctx.GetString("uid")})
}

func main() {
	engine := gin.Default()
	engine.LoadHTMLFiles("account.html", "manage.html")
	engine.GET("/login", Login)
	engine.GET("/account", Auth, Account)
	engine.GET("/manage", Auth, Manage)
	engine.Run("localhost:5678")
}

// cd jwt
// go run .
