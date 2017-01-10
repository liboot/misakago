package misakago

/*
  @brief http 请求封装
  @author xiyanxiyan10
  @data 2016/12/21
*/

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

//方法枚举值
const (
	GETMETHOD = iota
	POSTMETHOD
)

//http 请求句柄
type HttprequestHander struct {
	Url    string            //请求地址
	Body   []byte            //请求数据体
	Param  map[string]string //请求参数
	Method int               //请求使用的方法
}

//设置请求地址
func (handle *HttprequestHander) setUrl(url string) {
	handle.Url = url
}

//设置请求数据体
func (handle *HttprequestHander) setBody(data []byte) {
	handle.Body = data
}

//设置请求的方法
func (handle *HttprequestHander) setMethod(method int) {
	handle.Method = method
}

//设置请求参数
func (handle *HttprequestHander) setParam(param map[string]string) {
	for k, v := range param {
		handle.Param[k] = v
	}
}

//发送 post 请求
func (handle *HttprequestHander) post() ([]byte, error) {
	var dresult []byte
	body := bytes.NewBuffer([]byte(handle.Body))
	res, err := http.Post(handle.buildUrl(), "application/json;charset=utf-8", body)
	if err != nil {
		return dresult, nil
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return dresult, nil
	}
	return result, err
}

//发送 get 请求
func (handle *HttprequestHander) get() ([]byte, error) {
	var dresult []byte
	res, err := http.Get(handle.buildUrl())
	if err != nil {
		return dresult, nil
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return dresult, nil
	}
	return result, err
}

//构建使用的url地址，带参数
func (handle *HttprequestHander) buildUrl() string {
	u, _ := url.Parse(handle.Url)
	q := u.Query()
	for k, v := range handle.Param {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
