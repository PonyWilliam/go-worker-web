package handler

import (
	"context"
	"fmt"
	"github.com/PonyWilliam/go-WorkWeb/cache"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"strconv"
)

func AddWorker(c *gin.Context){
	user,_ := c.Get("username")
	if user != "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"请使用admin登陆",
		})
	}
	RspErr := gin.H{
		"code":500,
		"msg":"缺少参数",
	}
	temps := &works.Request_Workers{}
	i := 0

	if temp:=c.PostForm("name");temp!=""{
		temps.Name = temp
        i++
	}
	if temp:=c.PostForm("username");temp!=""{
		temps.Username = temp
        i++
	}
	if temp:=c.PostForm("password");temp!=""{
		temps.Password = temp
        i++
	}
	if temp:=c.PostForm("nums");temp!=""{
		temps.Nums = temp
        i++
	}
	if temp:=c.PostForm("place");temp!=""{
		temps.Place = temp
        i++
	}
	if temp:=c.PostForm("telephone");temp!=""{
		temps.Telephone = temp
        i++
	}
	if temp:=c.PostForm("score");temp!=""{
		new_temp,_ := strconv.ParseInt(temp,10,64)
		temps.Score = new_temp
        i++
	}
	if temp:=c.PostForm("description");temp!=""{
		temps.Description = temp
        i++
	}
	if temp:=c.PostForm("sex");temp!=""{
		temps.Sex = temp
        i++
	}
	if temp:=c.PostForm("mail");temp!=""{
		temps.Mail = temp
		i++
	}
	if temp:=c.PostForm("iswork");temp!=""{
		if temp == "true" || temp == "1"{
			temps.ISWork = true
		}else{
			temps.ISWork = false
		}
        i++
	}
	if temp:=c.PostForm("level");temp!=""{
		new_temp,_ := strconv.ParseInt(temp,10,64)
		temps.Level = new_temp
        i++
	}
	if i!=12 {
		c.JSON(200,RspErr)
		return
	}
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp,err := cl.CreateWorker(context.TODO(),temps)
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":err.Error(),
		})
		return
	}
	cache.DelCache("worker")
	c.JSON(200,gin.H{
		"code":200,
		"data":rsp.Id,
	})
}
func ChangeWorker(c *gin.Context){
	id := c.Param("id")
	if id == ""{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法识别到id",
		})
		return
	}
	new_id,_ := strconv.ParseInt(id,10,64)
	user,_ := c.Get("username")
	if user != "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"请使用admin登陆",
		})
	}
	RspErr := gin.H{
		"code":500,
		"msg":"缺少参数",
	}
	temps := &works.Request_Workers_Change{}
	i := 0
	//name: 123
	//level: 123
	//score: 100
	//sex: 男
	//telephone: 17794516068
	//place: 123
	//nums: 123
	//username: tjw666
	//password: tjw666
	//mail: tjw@dadiqq.cn
	//isWork: true
	//description: tjw666
	if temp:=c.PostForm("name");temp!=""{
		temps.Name = temp
		i++
	}
	if temp:=c.PostForm("username");temp!=""{
		temps.Username = temp
		i++
	}
	if temp:=c.PostForm("password");temp!=""{
		temps.Password = temp
		i++
	}
	if temp:=c.PostForm("nums");temp!=""{
		temps.Nums = temp
		i++
	}
	if temp:=c.PostForm("place");temp!=""{
		temps.Place = temp
		i++
	}
	if temp:=c.PostForm("telephone");temp!=""{
		temps.Telephone = temp
		i++
	}
	if temp:=c.PostForm("score");temp!=""{
		new_temp,_ := strconv.ParseInt(temp,10,64)
		temps.Score = new_temp
		i++
	}
	if temp:=c.PostForm("description");temp!=""{
		temps.Description = temp
		i++
	}
	if temp:=c.PostForm("sex");temp!=""{
		temps.Sex = temp
		i++
	}
	if temp:=c.PostForm("mail");temp!=""{
		temps.Mail = temp
		i++
	}
	if temp:=c.PostForm("iswork");temp!=""{
		if temp == "true" || temp == "1"{
			temps.ISWork = true
		}else{
			temps.ISWork = false
		}
		i++
	}
	if temp:=c.PostForm("level");temp!=""{
		new_temp,_ := strconv.ParseInt(temp,10,64)
		temps.Level = new_temp
		i++
	}
	if i!=12 {
		c.JSON(200,RspErr)
		return
	}
	temps.Id = new_id
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	_,err := cl.UpdateWorker(context.TODO(),temps)
	if err!=nil{
		c.JSON(200,err)
		return
	}
	cache.DelCache("worker")
	cache.DelCache(fmt.Sprintf("worker_%v",new_id))
	c.JSON(200,gin.H{
		"code":200,
		"msg":"操作成功",
	})
}
func DeleteWorker(c *gin.Context){
	user,ok := c.Get("username")
	if ok == false{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法验证你的身份信息",
		})
		return
	}
	if user!= "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权操作",
		})
		return
	}
	//通过delete方法删除
	id := c.Param("id")
	new_id,_ := strconv.ParseInt(id,10,64)
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp,err := cl.DeleteWorkerByID(context.TODO(),&works.Request_Workers_ID{Id: new_id})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":err.Error(),
		})
		return
	}
	cache.DelCache("worker")
	cache.DelCache(fmt.Sprintf("worker_%v",new_id))
	c.JSON(200,gin.H{
		"code":200,
		"msg":rsp.Message,
	})
}