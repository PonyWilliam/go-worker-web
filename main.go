package main

import (
	"WorkWeb/handler"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	"time"
)
func main() {
	consul:= consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1"}
	})
	router := gin.Default()
	service:= web.NewService(
		web.Name("go.micro.api.work"),
		web.Address("0.0.0.0:9999"),
		web.Registry(consul),
		web.Handler(router),
		)
	_ = service.Init()
	v1 := router.Group("work")
	v1.POST("/login/",handler.CheckPwd)
	v1.POST("/token/",handler.JWTAuthMiddleware(),handler.LoginByToken)
	v1.GET("/worker/:id",handler.JWTAuthMiddleware(),handler.GetUserInfoByID)
	v1.GET("/worker",handler.JWTAuthMiddleware(),handler.GetUserInfoAll)
	v1.GET("/workername/:name",handler.JWTAuthMiddleware(),handler.GetUserInfoByName)
	v1.GET("/workernum/:nums",handler.JWTAuthMiddleware(),handler.GetUserInfoByNums)
	v1.GET("/workerusername/:username",handler.JWTAuthMiddleware(),handler.GetUserInfoByUsername)
	v1.POST("/worker/borrow",handler.JWTAuthMiddleware(),handler.Borrow)
	v1.POST("/worker/other",handler.JWTAuthMiddleware(),handler.To_Other)
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AddAllowHeaders("token")
	corsConf.MaxAge = 3 * time.Hour
	router.Use(cors.New(corsConf))
	_ = router.Run()
	if err:=service.Run();err!=nil{
		fmt.Println(err.Error())
	}
}