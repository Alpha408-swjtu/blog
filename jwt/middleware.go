package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	if uid := getUidFromCookie1(ctx); uid != "" {
		ctx.Set("uid", uid)
		ctx.Next()
	} else {
		ctx.String(http.StatusForbidden, "请登录")
		ctx.Abort()
	}
}
