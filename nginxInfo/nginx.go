package nginxInfo

import (
	"encoding/json"
	"github.com/didi/nightingale/src/modules/collector/config"
	"github.com/didi/nightingale/src/modules/collector/util"
	"strconv"
	"strings"
	"time"
)

func NginxInfo() (string,error){
	res,err:= util.Get(config.GetYaml().NginxStatus, "")
	if err!=nil{
		config.Error.Println(err)
		return "接口查询失败", err
	}
	resNginxArr := strings.Split(res, "\n")
	//config.Info.Println(resNginxArr[2])
	count := strings.Split(resNginxArr[2], " ")
	successCount, err := strconv.Atoi(count[1])
	if err!=nil{
		config.Error.Println(err)
	}
	handsCount, err := strconv.Atoi(count[2])
	if err!=nil{
		config.Error.Println(err)
	}
	failCount := handsCount - successCount

	data := []config.Data{
		config.Data{
			Metric:    "nginx.query_count",
			Endpoint:  config.GetYaml().Ip,
			Tags:      "name=query_count",
			Value:     successCount,
			Timestamp: int(time.Now().Unix()),
			Step:      60,
		},
		config.Data{
			Metric:    "nginx.err_count",
			Endpoint:  config.GetYaml().Ip,
			Tags:      "name=err_count",
			Value:     failCount,
			Timestamp: int(time.Now().Unix()),
			Step:      60,
		},
	}
	jsonStr, err := json.MarshalIndent(data, "", " ")
	if err!=nil{
		config.Error.Println(err)
		return "fail",err
	}
	//config.Info.Println("formated: ", string(jsonStr))
	resData,err:= util.Post(config.GetYaml().Api, jsonStr, "application/json")
	return resData, err
}

