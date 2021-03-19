package handler

import (
	"context"
	"errors"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"time"
)

/*
	Login接口
*/
type Works struct{

}
var MySecret = []byte("rfiders") //密钥
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
func CheckPwd(c *gin.Context){
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp := &works.LoginResponse{}
	user := c.PostForm("user")
	rsp, err := cl.CheckSum(context.TODO(), &works.LoginRequest{User: user,
		Password: c.PostForm("password"),
	})
	if err!=nil{

		c.JSON(200,gin.H{
			"code":500,
			"msg":"验证密码出错",
		})
		return
	}
	if rsp.Code == false{
		c.JSON(200,gin.H{"code":500,"msg":"密码不正确"})
		return
	}
	token,_ := GenToken(user)
	c.JSON(200,gin.H{"code":200,"msg":"登陆成功","token":token})
}
func GenToken(username string)(string,error){
	c := MyClaims{username,jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		Issuer: "devicesManager",
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,c)
	return token.SignedString(MySecret)
}
func ParseToken(tokenString string)(*MyClaims,error){
	token,err := jwt.ParseWithClaims(tokenString,&MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return MySecret,nil
		},
	)
	if err!=nil{
		return nil,err
	}
	if claims,ok := token.Claims.(*MyClaims);ok && token.Valid{
		return claims,nil
	}
	return nil,errors.New("invalid token")
}
func LoginByToken(c *gin.Context){
	c.JSON(200,gin.H{
		"code":200,
		"msg":"success",
	})
}