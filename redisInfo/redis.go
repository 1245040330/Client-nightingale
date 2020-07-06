package redisInfo

import (
	"encoding/json"
	"github.com/didi/nightingale/src/modules/collector/config"
	"github.com/didi/nightingale/src/modules/collector/util"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
	"time"
)
//连接redis
func redisOpen() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetYaml().Redis.Addr,
		Password: config.GetYaml().Redis.Pass, // no password set
		DB:       0,                       // use default DB
	})
	return client
}

func RedisInfo() (string,error){
	data := []config.Data{}
	cache := redisOpen()
	//获取数据
	redisInfo := cache.Info()
	if strings.Index(redisInfo.String(), "ERR") != -1{
		config.Error.Println(redisInfo.String())
		return "", cache.Context().Err()
	}
	//连接总数 等待阻塞命令的客户端数  redis分配的内存总量  redis占用物理内存总量  redis内存使用最大值 每秒处理命令条数  拒绝连接数 rejected_connections
	arr := [...]string{"connected_clients", "blocked_clients", "used_memory", "used_memory_rss", "used_memory_peak", "instantaneous_ops_per_sec", "rejected_connections", "db0"}
	str := strings.Replace(redisInfo.String(), "\r", "", -1)
	dist1 := strings.Split(str, "\n")
	for _, key := range arr {
		for _, value := range dist1 {
			if strings.Index(value, key) != -1 {
				count := strings.Split(value, ":")
				if key == "db0" {
					numArr := strings.Split(count[1], ",")
					count[1] = strings.Split(numArr[0], "=")[1]
					key = key + ".num"
				}
				if strings.Index(value, "human") == -1 {
					intCount, err := strconv.Atoi(strings.Split(strings.Replace(count[1], "%", "", -1), ".")[0])
					if err!=nil{
						config.Error.Println(err)
					}
					data = append(data, config.Data{
						Metric:    "redis." + key,
						Endpoint:  config.GetYaml().Ip,
						Tags:      "name=" + key,
						Value:     intCount,
						Timestamp: int(time.Now().Unix()),
						Step:      60,
					})
				}
			}
		}
	}
	cache.Close()
	jsonStr, err := json.MarshalIndent(data, "", " ")
	if err!=nil{
		config.Error.Println(err)
		return "fail",err
	}
	//config.Info.Println("formated: ", string(jsonStr))
	res,err := util.Post(config.GetYaml().Api, jsonStr, "application/json")

	return res,err
}