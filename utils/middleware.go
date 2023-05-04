package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "*")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		context.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
		context.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
		context.Header("content-type", "application/json")
		context.Header("content-type", "text/html")
		// 设置返回格式是json
		//Release all option pre-requests
		if context.Request.Method == http.MethodOptions {
			context.JSON(http.StatusOK, "Options Request!")
		}
		context.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/" || strings.Contains(path, "/index") || strings.Contains(path, "/assets") {
			c.Next()
			return
		}
		token := c.Request.Header.Get("token")
		if path == "/ws" {
			query := c.Request.URL.Query()
			token = query["token"][0]
			appName := query["appName"][0]
			if appName == "" {
				c.Abort()
				c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "appName is not blank"})
				return
			}
		}
		if path == "/health" {
			query := c.Request.URL.Query()
			token = query["token"][0]
		}
		if !IsInvalidToken(token) {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "token invalid"})
			return
		} else {
			c.Next()
		}
	}
}

func WriteHtml(c *gin.Context) {
	accept := c.Request.Header.Get("Accept")
	flag := strings.Contains(accept, "text/html")
	if flag {
		content, err := os.ReadFile("dist/index.html")
		if (err) != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.Writer.WriteString("url:" + c.Request.URL.Path + " Not Found")
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Add("Content-Type", "text/html")
		c.Writer.Write(content)
		c.Writer.Flush()
	}
}
