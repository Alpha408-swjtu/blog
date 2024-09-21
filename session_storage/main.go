package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const COOKIE_NAME = "auth"

func Login(ctx *gin.Context) {
	//登录成功,给前端返回一个Set-Cookie
	uid := "8"
	payload := JwtPayload{UserDefined: map[string]any{"uid": uid}}
	token, _ := GenJWT(DefautHeader, payload, JWT_SECRET)
	fmt.Println("jwt token", token)

	ctx.JSON(http.StatusOK, gin.H{"uid": uid, "jwt": token})
}

func Account(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"uid": ctx.GetString("uid")})
}

func Manage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"uid": ctx.GetString("uid")})
}

func main() {
	engine := gin.Default()
	engine.LoadHTMLFiles("login.html", "account.html", "manage.html")
	engine.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", nil)
	})
	engine.GET("/login/body", Login)
	engine.GET("/account", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "account.html", nil)
	})
	engine.GET("/account/body", Auth, Account)
	engine.GET("/manage", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "manage.html", nil)
	})
	engine.GET("/manage/body", Auth, Manage)

	engine.Run("localhost:5678")
}

// cd session_storage
// go run .
