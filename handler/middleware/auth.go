package middleware

import (
	"blog/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UID_IN_TOKEN = "uid"
)

var (
	KeyConfig = util.CreateConfig("key")
)

// 从jwt里取出uid
func GetUidFromJwt(jwt string) int {
	_, payload, err := util.VerifyJwt(jwt, KeyConfig.GetString("jwt"))
	if err != nil {
		return 0
	}
	for k, v := range payload.UserDefined {
		if k == UID_IN_TOKEN {
			return int(v.(float64))
		}
	}
	return 0
}

// 从header里取出jwt，从而取出uid
func GetLoginUid(ctx *gin.Context) int {
	//方案一: 由客户端手动往request header里添加一个auth_token。redis访问次数少，在一浏览器窗口内切换用户登录不影响原来的窗口，缺点是refresh_token需要开启js访问权限，安全性略低。
	token := ctx.Request.Header.Get("auth_token")
	//方案二: 依靠浏览器自动回传的cookie，提取出refresh_token，再由服务端拿refresh_token查redis得到auth_token。增加了访问redis的频次，auth_token不需要传给前端；不同浏览器窗口共享同一个登录的uid。方案二实际上就是基于cookie的认证。
	// var token string
	//http协议里没有cookie这个概念，cookie本质上是header里的一对KV
	// for _, cookie := range strings.Split(ctx.Request.Header.Get("cookie"), ";") {
	// 	arr := strings.Split(cookie, "=")
	// 	key := strings.TrimSpace(arr[0])
	// 	value := strings.TrimSpace(arr[1])
	// 	if key == "refresh_token" {
	// 		token = database.GetToken(value)
	// 	}
	// }
	//或者直接使用封装好的Request.Cookies()
	// for _, cookie := range ctx.Request.Cookies() {
	// 	if cookie.Name == "refresh_token" {
	// 		fmt.Println(cookie.Value)
	// 		token = database.GetToken(cookie.Value)
	// 	}
	// }

	util.LogRus.Debugf("get token from header %s", token)
	return GetUidFromJwt(token)
}

// 身份认证中间件，先确保是登录状态
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		loginUid := GetLoginUid(ctx)
		if loginUid <= 0 {
			ctx.String(http.StatusForbidden, "auth failed") //返回403
			ctx.Abort()                                     //通过Abort()使中间件后面的handler不再执行，但是本handler还是会继续执行。所以下一行代码需要显式return
		} else {
			ctx.Set("uid", loginUid) //把登录的uid放入ctx中
			ctx.Next()
		}
	}
}
