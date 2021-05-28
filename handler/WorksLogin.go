package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PonyWilliam/go-WorkWeb/cache"
	borrow "github.com/PonyWilliam/go-borrow/proto"
	"github.com/PonyWilliam/go-common"
	product "github.com/PonyWilliam/go-product/proto"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/micro/go-micro/v2/client"
	"strings"
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
func LoginByToken(c *gin.Context){
	c.JSON(200,gin.H{
		"code":200,
		"msg":"身份有效",
	})
}
func CheckPwd(c *gin.Context){
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp := &works.LoginResponse{}
	username := c.PostForm("username")
	password := c.PostForm("password")
	rsp, err := cl.CheckSum(context.TODO(), &works.LoginRequest{User: username,
		Password: password,
	})
	fmt.Println(username,password)
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
	token,_ := GenToken(username)
	c.JSON(200,gin.H{"code":200,"msg":"登陆成功","token":token})
}
func GetQrcodeToken(c *gin.Context){
	user,ok := c.Get("username")
	if !ok{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法获取登录信息",
		})
	}
	token,_ := GenToken2(Strval(user))
	c.JSON(200,gin.H{
		"code":200,
		"token":token,
	})
}
func GetQrcodeInfo(c *gin.Context){
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	user,_ := c.Get("username")
	username := Strval(user)
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
func CheckQrcodeToken(c *gin.Context){
	//该接口就牛逼了，用于开门禁
	rfid := c.PostForm("rfid")
	rfids := strings.Split(rfid,",")
	if len(rfids) < 1{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"请携带rfid访问",
		})
	}
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	cl2 := product.NewProductService("go.micro.service.product",client.DefaultClient)
	user,_ := c.Get("username")
	username := Strval(user)
	res,err := cache.GetGlobalCache(fmt.Sprintf("worker_user_%v",username))
	rsp := &works.Response_Worker_Show{Worker: nil}
	if err != nil || err == redis.Nil{
		rsp,err := cl.FindWorkerByUserName(context.TODO(),&works.Request_Worker_User{Username: username})
		if err!=nil{
			c.JSON(200,gin.H{
				"code":500,
				"msg":"无法查询",
			})
			return
		}
		_ = cache.SetGlobalCache(fmt.Sprintf("worker_user_%v",username),rsp)
	}
	err = json.Unmarshal([]byte(res), &rsp)
	if err != nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"读取worker表失败",
		})
		return
	}
	//获取到了信息
	var products []*product.Response_ProductInfo
	for _,v := range rfids{
		res2,err := cl2.FindProductByRFID(context.TODO(),&product.Request_ProductRFID{Rfid: v})
		if err!=nil{
			c.JSON(200,gin.H{
				"code":500,
				"msg":err.Error(),
			})
			return
		}
		if res2.ProductIs == false{
			c.JSON(200,gin.H{
				"code":500,
				"msg":fmt.Sprintf("您租借的:%v不在库",res2.ProductName),
			})
			return
		}
		products = append(products,res2)
		//要求借出
	}
	for _,v := range products{
		BorrowBy(rsp,v)
	}
	c.JSON(200,gin.H{
		"code":200,
		"msg":"租借完成",
	})
}

func BorrowBy(worker *works.Response_Worker_Show,Product *product.Response_ProductInfo)(code int64,msg string){
	var err error
	cl1 := product.NewProductService("go.micro.service.product",client.DefaultClient)
	cl2 := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	borrowRequest := borrow.Borrow_Request{}
	fmt.Println(Product)
	if Product.ProductIs == false{
		code = 500
		msg = "商品不再库"
		return
	}
	if worker.Worker.Level < Product.ProductLevel{
		code = 500
		msg = "您无权借走该物品"
		return
	}
	ScheduleTime := time.Now().Unix() + 3600 * 24 * 3
	if Product.IsImportant{
		ScheduleTime = time.Now().Unix() + 3600 * 24 * 1
	}
	borrowRequest.ScheduleTime = ScheduleTime
	borrowRequest.WorkerId = worker.Worker.ID
	borrowRequest.ProductId = Product.Id
	rsp, _ := cl2.Borrow(context.TODO(), &borrowRequest)
	if rsp == nil || rsp.Status == 0{
		code = 500
		msg = "出借出错"
		return
	}
	Product.ProductIs = false
	Product.ProductBelongCustom = worker.Worker.ID
	temp2 := &product.Request_ProductInfo{}
	_ = common.SwapTo(Product, temp2)
	_, err = cl1.ChangeProduct(context.TODO(), temp2)
	if err!=nil{
		code = 500
		msg = err.Error()
	}
	code = 200
	msg = rsp.Message
	return
}


func GenToken(username string)(string,error){
	c := MyClaims{username,jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
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
//下面2种方式用于生成二维码token
func GenToken2(username string)(string,error){
	c := MyClaims{username,jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Second * 130).Unix(),
		Issuer: "qrcode",
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,c)
	return token.SignedString(MySecret)
}
func ParseToken2(tokenString string)(*MyClaims,error){
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
