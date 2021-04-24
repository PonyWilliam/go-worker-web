package main

import (
	"fmt"
	"github.com/PonyWilliam/go-WorkWeb/global"
	"github.com/PonyWilliam/go-WorkWeb/handler"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	"net/http"
)
func main() {
	consul:= consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"106.13.132.160"}
	})
	router := gin.Default()
	service:= web.NewService(
		web.Name("go.micro.api.work"),
		web.Address(":9999"),
		web.Registry(consul),
		web.Handler(router),
		)
	_ = service.Init()

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func(){
		err := http.ListenAndServe("9091",hystrixStreamHandler)
		if err != nil{
			log.Fatal(err)
		}
	}()
	v1 := router.Group("work")
	v1.POST("/login",handler.CheckPwd)
	v1.POST("/token/",handler.JWTAuthMiddleware(),handler.LoginByToken)
	v1.GET("/worker/:id",handler.GetUserInfoByID)
	v1.DELETE("/worker/:id",handler.JWTAuthMiddleware(),handler.DeleteWorker)
	v1.GET("/workers",handler.JWTAuthMiddleware(),handler.GetUserInfoAll)
	v1.POST("/worker",handler.JWTAuthMiddleware(),handler.AddWorker)
	v1.PUT("/worker/:id",handler.JWTAuthMiddleware(),handler.ChangeWorker)
	v1.GET("/workername/:name",handler.JWTAuthMiddleware(),handler.GetUserInfoByName)
	v1.GET("/workernum/:nums",handler.JWTAuthMiddleware(),handler.GetUserInfoByNums)
	v1.GET("/workerusername/:username",handler.JWTAuthMiddleware(),handler.GetUserInfoByUsername)
	v1.POST("/borrow/:pid",handler.JWTAuthMiddleware(),handler.Borrow)
	v1.GET("/borrow/user/",handler.JWTAuthMiddleware(),handler.GetBorrowByWorkerID)
	v1.POST("/return/:pid",handler.JWTAuthMiddleware(),handler.Return)
	v1.POST("/other",handler.JWTAuthMiddleware(),handler.To_Other)
	v1.POST("/confirm/:id",handler.JWTAuthMiddleware(),handler.Confirm)
	v1.POST("/reject/:id",handler.JWTAuthMiddleware(),handler.Reject)
	v1.GET("/req",handler.JWTAuthMiddleware(),handler.GetLogByReqWID)
	v1.GET("/rsp",handler.JWTAuthMiddleware(),handler.GetLogByRspWID)
	v1.GET("/",handler.JWTAuthMiddleware(),handler.GetAliCloudTempKey)
	router.Use(Cors())
	_ = router.Run()
	err := global.SetupRedisDb()
	if err != nil{
		log.Fatal(err)
	}
	if err:=service.Run();err!=nil{
		fmt.Println(err.Error())
	}
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, token, Token")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		fmt.Println("test")
		c.Next()
	}
}