package handler
// borrow 以及 return接口只有机器能访问到，会存储一个授权码到数据库或者通过其它等效方式
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PonyWilliam/go-WorkWeb/cache"
	borrowlog "github.com/PonyWilliam/go-borrow-logs/proto"
	borrow "github.com/PonyWilliam/go-borrow/proto"
	"github.com/PonyWilliam/go-common"
	product "github.com/PonyWilliam/go-product/proto"
	works "github.com/PonyWilliam/go-works/proto"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"strconv"
	"time"
)
func GetBorrowByWorkerID(c *gin.Context){
	user,ok := c.Get("username")
	now := c.Query("now")
	if !ok{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"token内容获取失败，请重新登陆",
		})
		return
	}
	user2 := Strval(user)
	worker,err := getUserInfo(user2)
	if worker == nil || worker.Worker == nil || err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"没有查询到你的信息",
		})
		return
	}
	fmt.Println(worker.Worker.ID)
	cl := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	rsp,err := cl.FindBorrowByWID(context.TODO(), &borrow.WID_Request{WID: worker.Worker.ID})
	if err!=nil || rsp==nil{
		fmt.Println(err)
		fmt.Println(rsp)
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法查询到信息",
		})
		return
	}
	if now == "now"{
		//只显示未归还的
		temp := &borrow.Borrowlogs_Response{}
			for _,v := range rsp.Logs{
			if v.ReturnTime == 0{
				temp.Logs = append(temp.Logs,v)
			}
		}
		c.JSON(200,temp.Logs)
		return
	}
	c.JSON(200,gin.H{
		"code":200,
		"data":rsp,
	})
}
func Borrow(c *gin.Context){
	/*
	(req.WorkerId,req.ProductId,req.ScheduleTime)
	*/
	username,_ := c.Get("username")
	if username!="admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权访问",
		})
		return
	}

	PID := c.Param("pid")
	if PID == ""{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"参数不全",
		})
		return
	}
	newPID,err := strconv.ParseInt(PID,10,64)
	cl1 := product.NewProductService("go.micro.service.product",client.DefaultClient)
	cl2 := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	borrowRequest := borrow.Borrow_Request{}

	user := c.PostForm("wid")
	new_user,_ := strconv.ParseInt(user,10,64)
	worker,err := getUserInfoByID(new_user)
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法获取用户信息",
		})
	}
	Product,err := cl1.FindProductByID(context.TODO(),&product.Request_ProductID{Id: newPID})
	fmt.Println(Product)
	if Product == nil || err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法获取商品信息",
		})
		return
	}
	if Product.ProductIs == false{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"商品不在库",
		})
		return
	}
	if worker.Worker.Level < Product.ProductLevel{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权限借走该商品",
		})
		return
	}
	ScheduleTime := time.Now().Unix() + 3600 * 24 * 3
	if Product.IsImportant{
		ScheduleTime = time.Now().Unix() + 3600 * 24 * 1
	}
	borrowRequest.ScheduleTime = ScheduleTime
	borrowRequest.WorkerId = worker.Worker.ID
	borrowRequest.ProductId = newPID
	rsp, _ := cl2.Borrow(context.TODO(), &borrowRequest)
	if rsp == nil || rsp.Status == 0{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"出借错误",
		})
		return
	}
	Product.ProductIs = false
	Product.ProductBelongCustom = worker.Worker.ID
	temp2 := &product.Request_ProductInfo{}
	_ = common.SwapTo(Product, temp2)
	_, err = cl1.ChangeProduct(context.TODO(), temp2)
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
func Return(c *gin.Context){
	pid := c.Param("pid")
	NewPid,_ := strconv.ParseInt(pid,10,64)
	user,_ := c.Get("username")
	if user != "admin"{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无权访问",
		})
		return
	}
	cl1 := product.NewProductService("go.micro.service.product",client.DefaultClient)
	cl2 := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	rsp,err := cl2.FindBorrowByPID(context.TODO(),&borrow.PID_Request{PID: NewPid})
	if err!=nil || rsp==nil || rsp.Logs==nil {
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法读取租借信息",
		})
		return
	}
	productInfo,err := cl1.FindProductByID(context.TODO(),&product.Request_ProductID{Id: NewPid})
	if err!=nil || productInfo == nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法读取产品信息",
		})
		return
	}
	if productInfo.ProductIs == true{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"当前商品已在库",
		})
		return
	}
	for _,v := range rsp.Logs{
		if v.ReturnTime == 0{
			_, err = cl2.Return(context.TODO(), &borrow.Returns_Request{Id: v.ID})
			if err!=nil{
				c.JSON(200,gin.H{
					"code":500,
					"msg":"归还时出错",
				})
				return
			}
			cache.DelCache(fmt.Sprintf("borrow_%v",v.ID))
			break
		}
	}
	productInfo.ProductIs = true
	productInfo.ProductBelongCustom = 0
	productInfo.ProductLocation = "库房"
	productRequest := &product.Request_ProductInfo{}
	_ = common.SwapTo(productInfo, productRequest)
	_,err = cl1.ChangeProduct(context.TODO(),productRequest)
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"更改产品状态出错",
		})
		return
	}
	cache.DelCache("borrow")
	c.JSON(200,gin.H{
		"code":200,
		"msg":"归还成功",
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
	cl3 := works.NewWorksService("go.micro.service.works",client.DefaultClient)
	rsp3, err := cl3.FindWorkerByNums(context.TODO(),&works.Request_Workers_Nums{Nums: new_other_id})
	if err!=nil || rsp3==nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"读取对方信息失败",
		})
		return
	}
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
			"msg":"越权访问，请重新登陆",
		})
		return
	}
	rsp, _ := cl.CheckToOther(context.TODO(), &borrow.ToOtherRequest{Id: new_log_id,Wid: new_other_id})
	if rsp.Status == 0{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"错误的wid",
		})
		return
	}
	rsp2,err := cl.FindBorrowByID(context.TODO(),&borrow.ID_Request{Id: new_log_id})
	if err != nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"findborrow错误",
		})
		return
	}
	t := &borrow.Borrowlog_Response{}
	if rsp2 == nil || rsp2== t{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"查询信息失败",
		})
		return
	}
	//等待他人确认后即可调用borrow服务改变状态
	cl2 := borrowlog.NewBorrowLogsService("go.micro.service.borrowlog",client.DefaultClient)
	rsp_temp,err := cl2.FindByLogID(context.TODO(),&borrowlog.Req_LogID{Logid: rsp2.ID})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":200,
			"msg":"查询logid失败",
		})
		return
	}
	test := 0
	if rsp_temp.Logs == nil{
		test = 1
	}else{
		for _,v := range rsp_temp.Logs{
			if v.Logid == new_log_id && v.Confirm == 3{
				//允许再次租借
				test = 1
			}
		}
	}
	if test == 0{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"有正在申请转借的记录，暂时无法重复申请",
		})
		return
	}
	rsp4,_ := cl2.ToOther(context.TODO(),&borrowlog.ReqToOther{ReqWID: worker.Worker.ID,RspWID: new_other_id,PID: rsp2.PID,Reason: Reason,Logid: rsp2.ID})
	if rsp4.Status == false{
		c.JSON(200,gin.H{
			"code":500,
			"msg":rsp4.Message,
		})
		return
	}
	cache.DelCache("borrow")
	cache.DelCache("borrow_log")
	c.JSON(200,gin.H{
		"code":200,
		"msg":"success",
	})
}
func Confirm(c *gin.Context){
	id := c.Param("id")
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
	rsp,err := cl.FindByID(context.TODO(),&borrowlog.Req_Id{Id: new_id})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":200,
			"msg":"无法读取对应id",
		})
		return
	}
	if rsp.Confirm != 1 {
		c.JSON(200,gin.H{
			"code":200,
			"msg":"您已处理过",
		})
		return
	}
	if rsp.RspWID != worker.Worker.ID{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权确认",
		})
		return
	}
	//权限判断完毕
	_, err = cl.Confirm(context.TODO(), &borrowlog.Req_Confirm{ID: new_id})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":err.Error(),
		})
		return
	}
	cl2 := borrow.NewBorrowService("go.micro.service.borrow",client.DefaultClient)
	rsp2,err := cl2.ToOther(context.TODO(),&borrow.ToOtherRequest{Id: rsp.Logid,Wid: rsp.RspWID})
	if err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":err.Error(),
		})
		return
	}
	//改变产品所有者
	cl3 := product.NewProductService("go.micro.service.product",client.DefaultClient)
	products,err := cl3.FindProductByID(context.TODO(),&product.Request_ProductID{Id: rsp.PID})
	if err != nil{
		c.JSON(200,"fail")
		return
	}
	_,_ = cl3.ChangeProduct(context.TODO(),&product.Request_ProductInfo{Id: rsp.PID,ProductBelongCustom: rsp.RspWID,ProductName: products.ProductName,ProductDescription: products.ProductDescription,ProductLevel: products.ProductLevel,ProductBelongCategory: products.ProductBelongCategory,ProductBelongArea: products.ProductBelongArea,ProductIs: products.ProductIs,ProductRfid: products.ProductRfid,ProductLocation: products.ProductLocation,ImageId: products.ImageId,IsImportant: products.IsImportant})
	//这TM就是给自己挖的大坑。直到现在才被发现。
	c.JSON(200,rsp2.Message)
}
func Reject(c *gin.Context){
	id := c.Param("id")
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
	if rsp == nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"无法读取对应信息",
		})
		return
	}
	if rsp.Confirm != 1 {
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您已处理过",
		})
		return
	}
	if rsp.RspWID != worker.Worker.ID{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"您无权确认",
		})
		return
	}
	rsp2,_ := cl.Reject(context.TODO(),&borrowlog.Req_Reject{ID: new_id})
	if rsp2.Status == false{
		c.JSON(200,gin.H{
			"code":500,
			"msg":rsp2.Message,
		})
		return
	}

	c.JSON(200,gin.H{
		"code":200,
		"msg":rsp2.Message,
	})
}
func GetLogByReqWID(c *gin.Context){
	user,ok := c.Get("username")
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
	rsp, _ := cl.FindByReqWID(context.TODO(),&borrowlog.Req_Wid{Wid: worker.Worker.ID})
	c.JSON(200,gin.H{
		"code":200,
		"msg":rsp.Logs,
	})
}
func GetLogByRspWID(c *gin.Context){
	user,ok := c.Get("username")
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
	rsp, _ := cl.FindByRspWID(context.TODO(),&borrowlog.Req_Wid{Wid: worker.Worker.ID})
	c.JSON(200,gin.H{
		"code":200,
		"data":rsp.Logs,
	})
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