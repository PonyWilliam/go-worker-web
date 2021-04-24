package handler

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts2 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
)
//request.DurationSeconds = common.Uint64Ptr(7200)
func GetTencentCloudTempApi(c *gin.Context){
	credential := common.NewCredential(
		"AKIDeZyLcc6sESVv73DHkGakU0KOCl5CR7qM",
		"SghXv53QWO0Wiky1rRq0wf8SWdT2MPDk",
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sts.tencentcloudapi.com"
	client, _ := sts2.NewClient(credential, "ap-chengdu", cpf)

	request := sts2.NewGetFederationTokenRequest()

	request.Name = common.StringPtr("test")
	policy := string(`{
"version": "2.0",
"statement": [
  {
    "action": [
      "name/cos:PutObject"
    ],
    "effect": "allow",
    "resource": [
      "qcs::cos:ap-chengdu:uid/1257689370:arcsoft-1257689370/*"
    ]
  }
]
}`)
	request.Policy = common.StringPtr(policy)
	response, err := client.GetFederationToken(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		c.Status(500)
		return
	}
	if err != nil {
		panic(err)
	}
	c.JSON(200,response)
}
func GetAliCloudTempKey(c *gin.Context){
	client, err := sts.NewClientWithAccessKey("cn-guangzhou", "LTAI5tHp7SKmL9h1eTnu7C72", "5tNzDpvZBH9DAoD9rFJQBBHHTkyQKY")
	if client == nil{
		c.JSON(200,gin.H{
			"code":500,
			"msg":"未知错误",
		})
		return
	}
	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//设置参数。关于参数含义和设置方法，请参见API参考。
	request.RoleArn = "acs:ram::1965385939128175:role/arcsoft"
	request.RoleSessionName = "myRole"

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		c.JSON(200,gin.H{
			"code":500,
			"msg":err.Error(),
		})
	}
	c.JSON(200,response)
}