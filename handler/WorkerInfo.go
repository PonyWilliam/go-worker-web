package handler

import (
	"context"
	"fmt"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"strconv"
)

func GetUserInfoByID(c *gin.Context){
	//获取基本信息，渲染到前端
	user,_ := c.Get("username")
	if user != "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权访问",
		})
	}
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	id := c.Param("id")
	rsp := &works.Response_Worker_Show{Worker: nil}
	newid, _ := strconv.ParseInt(id, 10, 64)
	rsp,err := cl.FindWorkerByID(context.TODO(),&works.Request_Workers_ID{Id: newid})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法查询",
		})
	}
	c.JSON(200,rsp)
}
func GetUserInfoAll(c *gin.Context){
	user,ok := c.Get("username")
	fmt.Println("user")
	fmt.Println(user)
	if ok == false{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"can not read id",
		})
		return
	}
	if user != "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"请使用管理员账号查询",
		})
		return
	}
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp := &works.Response_Workers_Show{Workers: nil}
	rsp,err := cl.FindAll(context.TODO(),&works.Request_Null{})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法查询",
		})
		return
	}
	c.JSON(200,rsp)
}
func GetUserInfoByNums(c *gin.Context){
	user,_ := c.Get("username")
	if user != "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权访问",
		})
	}
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	nums := c.Param("nums")
	rsp := &works.Response_Worker_Show{Worker: nil}
	newnums, _ := strconv.ParseInt(nums, 10, 64)
	rsp,err := cl.FindWorkerByNums(context.TODO(),&works.Request_Workers_Nums{Nums: newnums})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法查询",
		})
		return
	}
	c.JSON(200,rsp)
}
func GetUserInfoByName(c *gin.Context){
	//获取基本信息，渲染到前端
	user,_ := c.Get("username")
	if user != "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权访问",
		})
	}
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	name := c.Param("name")
	rsp := &works.Response_Workers_Show{Workers: nil}
	rsp,err := cl.FindWorkerByName(context.TODO(),&works.Request_Workers_Name{Name: name})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法查询",
		})
		return
	}
	c.JSON(200,rsp)
}

func GetUserInfoByUsername(c *gin.Context){
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	username := c.Param("username")
	username2,_ := c.Get("username")
	if username != username2 && username2 != "admin" {
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权访问",
		})
		return
	}
	rsp := &works.Response_Worker_Show{Worker: nil}
	rsp,err := cl.FindWorkerByUserName(context.TODO(),&works.Request_Worker_User{Username: username})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法查询",
		})
		return
	}
	c.JSON(200,rsp)
}