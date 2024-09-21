package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getUidFromHeader(ctx *gin.Context) string {
	//从request header里取得jwt
	if token := strings.TrimSpace(ctx.Request.Header.Get("jwt")); token != "" {
		if _, payload, err := VerifyJwt(token, JWT_SECRET); err == nil {
			if uid, exists := payload.UserDefined["uid"]; exists {
				return uid.(string)
			}
		}
	}
	return ""
}

func Auth(ctx *gin.Context) {
	if uid := getUidFromHeader(ctx); uid != "" {
		ctx.Set("uid", uid)
		ctx.Next()
	} else {
		ctx.String(http.StatusForbidden, "请登录")
		ctx.Abort()
	}
}
