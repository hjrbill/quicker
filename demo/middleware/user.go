package middleware

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

func GetUserInfo(ctx *gin.Context) {
	userName, err := url.QueryUnescape(ctx.Request.Header.Get("UserName")) //从 request header 里获得 UserName
	if err == nil {
		ctx.Set("user_name", userName) //把 UserName 放到 gin.Context 里
	}
}
