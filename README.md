# openapiLib

本地下载代码运行  
go get github.com/wf42988/openapiLib  
go.mod 中 replace openapiLib => ../openapiLib  
import "openapiLib"

git下载运行  
go.mod 中  
require github.com/wf42988/openapiLib latest  
import github.com/wf42988/openapiLib"   
```golang
// openapiTest project main.go
package main

import (
    "encoding/json"
    "fmt"
	"openapiLib"
	"strconv"
	"time"
)

func strToLong(timeStr string) int64 {
	timeTemplate := "2006-01-02 15:04:05"
	stamp, _ := time.ParseInLocation(timeTemplate, timeStr, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	return stamp.Unix() * 1000
}
func main() {
	// openApiGatewayClient := openapiLib.OpenApiGatewayClient{}
	// openApiGatewayClient.ApiServerPath = "http://z.smarteoc.com:8897/"
	// openApiGatewayClient.Version = "3.2.0"
	// openApiGatewayClient.ClientId = "5zsf8f9332a8fca1"
	// openApiGatewayClient.AccessKeyId = "5zrf8f9333a7cf01"
	// openApiGatewayClient.AccessKeySecret = "ZGVlNDYzODA3MGUwNDg2MGFlZDM4NTFmMDU5Yzg0ZGI="
	// token_ret := openApiGatewayClient.GetTokenClientCredentials()
	// fmt.Println(token_ret)

	openApiGatewayClient := openapiLib.OpenApiGatewayClient{ApiServerPath: "http://z.smarteoc.com:8897/", Version: "3.2.0", ClientId: "5zsf8f9332a8fca1", AccessKeyId: "5zrf8f9333a7cf01", AccessKeySecret: "ZGVlNDYzODA3MGUwNDg2MGFlZDM4NTFmMDU5Yzg0ZGI="}
	// token_ret := openApiGatewayClient.GetTokenClientCredentials()
	// fmt.Println(token_ret)
	// var token_ret_dic map[string]interface{}
	// json.Unmarshal([]byte(token_ret), &token_ret_dic)
	// access_token, ok := token_ret_dic["access_token"]
	// if ok == false {
	// 	return
	// }
	// fmt.Println(access_token)

	token_ret := openApiGatewayClient.GetTokenPassword("pttest", "123456")
	fmt.Println(token_ret)
	var token_ret_dic map[string]interface{}
	json.Unmarshal([]byte(token_ret), &token_ret_dic)
	access_token, ok := token_ret_dic["access_token"]
	if ok == false {
		return
	}
	// refresh_token, ok := token_ret_dic["refresh_token"]
	// if ok == false {
	// 	return
	// }

	// refresh_token_ret := openApiGatewayClient.RefreshToken(refresh_token.(string))
	// fmt.Println(refresh_token_ret)
	// var refresh_token_ret_dic map[string]interface{}
	// json.Unmarshal([]byte(refresh_token_ret), &refresh_token_ret_dic)
	// refresh_access_token, ok := refresh_token_ret_dic["access_token"]
	// if ok == false {
	// 	return
	// }
	// access_token = refresh_access_token

	dev_url := "http://z.smarteoc.com:8897/"
	// test_url := "http://10.1.170.30:80/"
	// product_url := "http://openapi3.smarteoc.com/"
	// pre_url := "http://openapi3.pre.smarteoc.com/"
	// k8s_url := "http://10.1.170.61:30100/"
	url := dev_url
	start := strToLong("2019-11-19 00:00:00")
	end := strToLong("2019-11-20 00:00:00")
	startStr := strconv.FormatInt(start, 10)
	endStr := strconv.FormatInt(end, 10)
	para_dic := make(map[string]string)
	para_dic["pointIds"] = "169913186689007616"
	para_dic["startTime"] = startStr
	para_dic["endTime"] = endStr
	// ret := openApiGatewayClient.DoAction(url, "metricDatum", "queryRawMetricDatum", para_dic, access_token.(string))
	ret := openApiGatewayClient.DoActionWithRequestId("requestId111", url, "metricDatum", "queryRawMetricDatum", para_dic, access_token.(string))
	fmt.Println(ret)
}
```