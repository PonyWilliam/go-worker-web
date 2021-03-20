package handler
// borrow 以及 return接口只有机器能访问到，会存储一个授权码到数据库或者通过其它等效方式
import (
	"context"
	"encoding/json"
	"fmt"
	borrowlog "github.com/PonyWilliam/go-borrow-logs/proto"
	borrow "github.com/PonyWilliam/go-borrow/proto"
	product "github.com/PonyWilliam/go-product/proto"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"strconv"
	"time"
)
func Borrow(c *gin.Context){
	/*
	(req.WorkerId,req.ProductId,req.ScheduleTime)
	*/
	PID := c.PostForm("PID")
	if PID == ""{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"参数不全",
		})
		return
	}
	ScheduleTime := time.Now().Unix() + 3600 * 24 * 7
	newPID,err := strconv.ParseInt(PID,10,64)
	cl1 := product.NewProductService("go.micro.service.product",client.DefaultClient)
	cl2 := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	borrowRequest := borrow.Borrow_Request{}

	user,ok := c.Get("username")
	newuser := Strval(user)
	if !ok{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"token内容获取失败，请重新登陆",
		})
		return
	}
	worker,err := getUserInfo(newuser)
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法获取用户信息",
		})
	}
	Product,err := cl1.FindProductByID(context.TODO(),&product.Request_ProductID{Id: newPID})
	if Product == nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法获取商品信息",
		})
		return
	}
	if worker.Worker.Level < Product.Info.ProductLevel{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权限借走该商品",
		})
		return
	}
	borrowRequest.ScheduleTime = ScheduleTime
	borrowRequest.WorkerId = worker.Worker.ID
	borrowRequest.ProductId = newPID
	rsp, _ := cl2.Borrow(context.TODO(), &borrowRequest)
	if rsp.Status == 0{
		c.JSON(200,gin.H{
			"code":500,
			"msg":rsp.Message,
		})
		return
	}
	Product.Info.ProductIs = false
	_, err = cl1.ChangeProduct(context.TODO(), Product.Info)
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":err.Error(),
		})
	}
	c.JSON(200,gin.H{
		"code":200,
		"msg":rsp.Message,
	})
}
func To_Other(c *gin.Context){
	log_id := c.PostForm("id")//borrow表内ID
	Reason := c.PostForm("reason")
	other_id := c.PostForm("wid")//其他人的WID
	if log_id == "" || Reason == "" || other_id == ""{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"参数不全",
		})
		return
	}
	new_log_id,err := strconv.ParseInt(log_id,10,64)
	new_other_id,err := strconv.ParseInt(other_id,10,64)
	cl := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	user,ok := c.Get("username")
	newuser := Strval(user)
	if !ok{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"token内容获取失败，请重新登陆",
		})
		return
	}
	worker,err := getUserInfo(newuser)
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法获取用户信息",
		})
		return
	}
	if worker == nil {
		c.JSON(200,gin.H{
			"code":500,
			"msg":"用户表发生未知错误，请联系管理员",
		})
		return
	}
	rsp, _ := cl.CheckToOther(context.TODO(), &borrow.ToOtherRequest{Id: new_log_id,Wid: new_other_id})
	if rsp.Status == 0{
		fmt.Println("在1错误")
		c.JSON(200,gin.H{
			"code":500,
			"msg":"check错误",
		})
		return
	}
	rsp2,err := cl.FindBorrowByID(context.TODO(),&borrow.ID_Request{Id: new_log_id})
	if err != nil{
		fmt.Println("在2错误")
		c.JSON(200,gin.H{
			"code":500,
			"msg":"findborrow错误",
		})
		return
	}
	if rsp2 == nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"查询信息失败",
		})
		return
	}
	fmt.Println(rsp2)
	//等待他人确认后即可调用borrow服务改变状态
	cl2 := borrowlog.NewBorrowLogsService("go.micro.service.borrowlog",client.DefaultClient)
	_,err = cl2.ToOther(context.TODO(),&borrowlog.ReqToOther{ReqWID: worker.Worker.ID,RspWID: new_other_id,PID: rsp2.PID,Reason: Reason})
	if err!=nil{
		fmt.Println("在3错误")
		c.JSON(200,gin.H{
			"code":500,
			"msg":err.Error(),
		})
		return
	}
	fmt.Println("hello")
	c.JSON(200,gin.H{
		"code":200,
		"msg":"success",
	})
}
func Confirm(c *gin.Context){
	id := c.PostForm("id")
	user,ok := c.Get("username")
	new_id,err := strconv.ParseInt(id,10,64)
	if !ok{
		c.JSON(200, gin.H{
			"code":500,
			"msg":"非法访问",
		})
		return
	}
	new_user := Strval(user)
	worker,err := getUserInfo(new_user)
	if err != nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法读取到你的信息",
		})
		return
	}
	cl := borrowlog.NewBorrowLogsService("go.micro.service.borrowlog",client.DefaultClient)
	rsp, _ := cl.FindByID(context.TODO(), &borrowlog.Req_Id{Id: new_id})
	if rsp.RspWID != worker.Worker.ID{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权确认",
		})
		return
	}
	//权限判断完毕
	cl2 := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	cl2.ToOther(context.TODO(),&borrow.ToOtherRequest{})
}
//下面是小工具
func getUserInfo(username string)(*works.Response_Worker_Show,error){
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp, err := cl.FindWorkerByUserName(context.TODO(), &works.Request_Worker_User{Username: username})
	if err!=nil {
		return rsp,err
	}
	return rsp,nil
}
func getUserInfoByID(ID int64)(*works.Response_Worker_Show,error){
	cl := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp, err := cl.FindWorkerByID(context.TODO(), &works.Request_Workers_ID{Id: ID})
	if err!=nil {
		return rsp,err
	}
	return rsp,nil
}
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}