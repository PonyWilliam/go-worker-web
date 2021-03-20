package handler

import (
	"context"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"strconv"
)

func SetWorker(c *gin.Context){
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
	if temp:=c.PostForm("Level");temp!=""{
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
		c.JSON(200,err)
		return
	}
	c.JSON(200,rsp)
}
func ChangeWorker(c *gin.Context){
	id := c.Param("id")
	if id == ""{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法识别到id",
		})
	}
	new_id,_ := strconv.ParseInt(id,10,64)
	RspErr := gin.H{
		"code":500,
		"msg":"没有任何更改",
	}
	temps := &works.Request_Workers_Change{}
	temps.Id = new_id
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
	if temp:=c.PostForm("Level");temp!=""{
		new_temp,_ := strconv.ParseInt(temp,10,64)
		temps.Level = new_temp
		i++
	}
	if i==0 {
		c.JSON(200,RspErr)
		return
	}
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp,err := cl.UpdateWorker(context.TODO(),temps)
	if err!=nil{
		c.JSON(200,err)
		return
	}
	c.JSON(200,rsp)
}