package util

import (
	"github.com/didi/nightingale/src/modules/collector/config"
	"github.com/imroc/req"
	"time"
)



// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func Post(url string, data interface{}, contentType string) (string,error) {
	heade := req.Header{
		"Content-Type": contentType,
	}
	req.SetTimeout(2 * time.Second)
	r, err := req.Post(url, heade, data)

	if err!=nil{
		config.Error.Println(err)
		return "接口超时", err
	}
	//Info.Println(r.ToString())
	res, err := r.ToString()
	if err!=nil{
		config.Error.Println(err)
	}
	return res, err
}

func Get(url string, data interface{}) (string, error) {
	req.SetTimeout(2 * time.Second)
	r, err := req.Get(url, data)
	if err!=nil{
		config.Error.Println(err)
		return "接口超时", err
	}
	//Info.Println(r.ToString())
	res, err := r.ToString()
	if err!=nil{
		config.Error.Println(err)
	}
	return res ,err
}