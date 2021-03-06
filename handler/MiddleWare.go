package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中,在跨域中也允许该字段跨域
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请携带token访问",
			})
			c.Abort()
			return
		}
		// 按空格分割

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(authHeader)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "token失效",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Set("id",mc.Id)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

func JWTAuthMiddleware2() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中,在跨域中也允许该字段跨域
		authHeader := c.Query("secret")//采用get明文传输
		if authHeader == "" {
			authHeader = c.PostForm("secret")
			if authHeader == ""{
				c.JSON(http.StatusOK, gin.H{
					"code": 2003,
					"msg":  "请携带token访问",
				})
				c.Abort()
				return
			}
		}
		// 按空格分割

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken2(authHeader)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "token失效",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}