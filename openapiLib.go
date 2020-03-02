// openapiLib project openapiLib.go
package openapiLib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type OpenApiGatewayClient struct {
	ApiServerPath   string
	Version         string
	ClientId        string
	AccessKeyId     string
	AccessKeySecret string
}

func signature(service string, para_dic map[string]string, ms int64, secret string) string {
	keys := make([]string, 0, len(para_dic))
	for key := range para_dic {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	url_str := ""
	for _, key := range keys {
		value := para_dic[key]
		url_str = url_str + "&" + url.QueryEscape(key) + "=" + url.QueryEscape(value)
	}
	url_str = url_str[1:]
	stringToSign := "POST&" + url.QueryEscape("/") + "&" + url.QueryEscape(url_str)
	final_secret := secret + "&"
	hmacsha256 := hmac.New(sha256.New, []byte(final_secret))
	io.WriteString(hmacsha256, stringToSign)
	signature_str := base64.StdEncoding.EncodeToString(hmacsha256.Sum(nil))
	return signature_str
}

func doAction(requestId string, url string, service string, method string, paraDic map[string]string, accessKeyId string, secret string, access_token string) string {
	ms := time.Now().Unix() * 1000
	// var ms int64 = 1574141867043
	url_final := url + service
	fina_para_dic := make(map[string]string)
	fina_para_dic["action"] = method
	fina_para_dic["version"] = "3.2.0"
	fina_para_dic["timestamp"] = strconv.FormatInt(ms, 10)
	fina_para_dic["signatureMethod"] = "HMAC-SHA256"
	fina_para_dic["signatureVersion"] = "1.0"
	fina_para_dic["accessKeyId"] = accessKeyId
	if requestId != "" {
		fina_para_dic["requestId"] = requestId
	}
	for key, value := range paraDic {
		fina_para_dic[key] = value
		// fmt.Printf("%s==%s\n", key, value)
	}
	signature_str := signature(service, fina_para_dic, ms, secret)
	fina_para_dic["signature"] = signature_str
	fina_para_dic["access_token"] = access_token
	var r http.Request
	r.ParseForm()
	for key, value := range fina_para_dic {
		r.Form.Add(key, value)
	}
	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest("POST", url_final, strings.NewReader(bodystr)) //application/x-www-form-urlencoded
	// request, err := http.NewRequest("POST", url_final+"/", strings.NewReader(string(bytesData)))
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Set("Content-Type", "application/json;charset=UTF-8") //;charset=UTF-8
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str
}

func getToken(tokern_ret string) string {
	var token_ret_dic map[string]interface{}
	json.Unmarshal([]byte(tokern_ret), &token_ret_dic)
	access_token, ok := token_ret_dic["access_token"]
	if ok == false {
		return tokern_ret
	}
	// fmt.Println(access_token)
	ret := access_token.(string)
	return ret
}

func (openApiGatewayClient OpenApiGatewayClient) DoAction(url string, service string, method string, paraDic map[string]string, access_token string) string {
	str := doAction("", url, service, method, paraDic, openApiGatewayClient.AccessKeyId, openApiGatewayClient.AccessKeySecret, access_token)
	return str
}

func (openApiGatewayClient OpenApiGatewayClient) DoActionWithRequestId(requestId string, url string, service string, method string, paraDic map[string]string, access_token string) string {
	str := doAction(requestId, url, service, method, paraDic, openApiGatewayClient.AccessKeyId, openApiGatewayClient.AccessKeySecret, access_token)
	return str
}

func (openApiGatewayClient OpenApiGatewayClient) GetTokenClientCredentials() string {
	var r http.Request
	url := openApiGatewayClient.ApiServerPath + "oauth/token"
	r.ParseForm()
	r.Form.Add("grant_type", "client_credentials")
	r.Form.Add("client_id", openApiGatewayClient.ClientId)
	r.Form.Add("client_secret", openApiGatewayClient.AccessKeyId)
	bodystr := strings.TrimSpace(r.Form.Encode())
	request, err := http.NewRequest("POST", url, strings.NewReader(bodystr)) //application/x-www-form-urlencoded
	// request, err := http.NewRequest("POST", url_final+"/", strings.NewReader(string(bytesData)))
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Set("Content-Type", "application/json;charset=UTF-8") //;charset=UTF-8
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str
}

func (openApiGatewayClient OpenApiGatewayClient) GetTokenPassword(username string, password string) string {
	var r http.Request
	r.ParseForm()
	r.Form.Add("grant_type", "password")
	r.Form.Add("client_id", openApiGatewayClient.ClientId)
	r.Form.Add("client_secret", openApiGatewayClient.AccessKeyId)
	r.Form.Add("username", username)
	r.Form.Add("password", password)
	bodystr := strings.TrimSpace(r.Form.Encode())
	url := openApiGatewayClient.ApiServerPath + "oauth/token"
	request, err := http.NewRequest("POST", url, strings.NewReader(bodystr)) //application/x-www-form-urlencoded
	// request, err := http.NewRequest("POST", url_final+"/", strings.NewReader(string(bytesData)))
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Set("Content-Type", "application/json;charset=UTF-8") //;charset=UTF-8
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str
}

func (openApiGatewayClient OpenApiGatewayClient) RefreshToken(refreshToken string) string {
	var r http.Request
	r.ParseForm()
	r.Form.Add("grant_type", "refresh_token")
	r.Form.Add("refresh_token", refreshToken)
	r.Form.Add("client_id", openApiGatewayClient.ClientId)
	r.Form.Add("client_secret", openApiGatewayClient.AccessKeyId)
	bodystr := strings.TrimSpace(r.Form.Encode())
	url := openApiGatewayClient.ApiServerPath + "oauth/token"
	request, err := http.NewRequest("POST", url, strings.NewReader(bodystr)) //application/x-www-form-urlencoded
	// request, err := http.NewRequest("POST", url_final+"/", strings.NewReader(string(bytesData)))
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Set("Content-Type", "application/json;charset=UTF-8") //;charset=UTF-8
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str
}
