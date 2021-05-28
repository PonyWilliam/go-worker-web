package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PonyWilliam/go-WorkWeb/cache"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/micro/go-micro/v2/client"
	"strconv"
)

func GetUserInfoByID(c *gin.Context){
	//获取基本信息，渲染到前端
	id := c.Param("id")
	newid, _ := strconv.ParseInt(id, 10, 64)
	res,err := cache.GetGlobalCache(fmt.Sprintf("worker_%v",newid))
	if err!= nil || err == redis.Nil{
		cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
		rsp := &works.Response_Worker_Show{Worker: nil}
		rsp,err := cl.FindWorkerByID(context.TODO(),&works.Request_Workers_ID{Id: newid})
		if err!=nil{
			c.JSON(200,gin.H{
				"code":500,
				"msg":"无法查询",
			})
			return
		}
		_ = cache.SetGlobalCache(fmt.Sprintf("worker_%v", newid), rsp)
		c.JSON(200,gin.H{
			"code":200,
			"data":rsp,
		})
		return
	}
	result := &works.Response_Worker_Show{}
	_ = json.Unmarshal([]byte(res), &result)
	c.JSON(200,gin.H{
		"code":200,
		"data":result,
	})
}
func GetUserInfoAll(c *gin.Context){
	_,ok := c.Get("username")
	if !ok{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"请携带正确的token访问",
		})
	}
	res,err := cache.GetGlobalCache("worker")
	if err != nil || err == redis.Nil{
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
		_ = cache.SetGlobalCache("worker",rsp)
		c.JSON(200,gin.H{
			"code":200,
			"data":rsp,
		})
		return
	}
	result := &works.Response_Workers_Show{}
	_ = json.Unmarshal([]byte(res), &result)
	c.JSON(200,gin.H{
		"code":200,
		"data":result,
	})
}
func GetUserInfoByNums(c *gin.Context){
	_,ok := c.Get("username")
	if !ok{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"请携带正确的token访问",
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
	c.JSON(200,gin.H{
		"code":200,
		"data":rsp,
	})
}
func GetUserInfoByName(c *gin.Context){
	//获取基本信息，渲染到前端
	_,ok := c.Get("username")
	if !ok{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"请携带正确的token访问",
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
	c.JSON(200,gin.H{
		"code":200,
		"data":rsp,
	})
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
	res,err := cache.GetGlobalCache(fmt.Sprintf("worker_user_%v",username))
	if err != nil || err == redis.Nil{
		rsp := &works.Response_Worker_Show{Worker: nil}
		rsp,err := cl.FindWorkerByUserName(context.TODO(),&works.Request_Worker_User{Username: username})
		if err!=nil{
			c.JSON(200,gin.H{
				"code":500,
				"msg":"无法查询",
			})
			return
		}
		_ = cache.SetGlobalCache(fmt.Sprintf("worker_user_%v",username),rsp)
		c.JSON(200,gin.H{
			"code":200,
			"data":rsp,
		})
		return
	}
	result := &works.Response_Worker_Show{}
	_ = json.Unmarshal([]byte(res), &result)
	c.JSON(200,gin.H{
		"code":200,
		"data":result,
	})
}