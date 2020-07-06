package pgsqlInfo

import (
	"database/sql"
	"encoding/json"
	"github.com/didi/nightingale/src/modules/collector/config"
	"github.com/didi/nightingale/src/modules/collector/util"
	"time"
	_ "github.com/lib/pq"
)

//pgsql数据采集
func SqlInfo() (string,error){
	db, err := sql.Open("postgres", "port="+config.GetYaml().Pgsql.Port+" user=postgres password="+config.GetYaml().Pgsql.Pass+" dbname="+config.GetYaml().Pgsql.Dbname+" sslmode=disable")
	if err!=nil{
		config.Error.Println(err)
	}
	config.Info.Println("pgsqlOpen")
	count, err:=sqlSelect(db)
	if err!=nil{
		config.Error.Println(err)
		return "fail",err
	}
	config.Info.Println("*sqlSelect")
	db.Close()
	config.Info.Println("*sqlClose")

	res,err:=pgsqlPostData(count)
	if err!=nil{
		config.Error.Println(err)
	}
	return res,err
}

//查询连接数
func sqlSelect(db *sql.DB) (int64,error){
	rows, err := db.Query("select count('count') from pg_stat_activity")
	if err!=nil{
		config.Error.Println(err)
		return 0,err
	}
	var count int64
	for rows.Next() {
		err = rows.Scan(&count)
	}
	return count,err
}

func pgsqlPostData(count int64) (string, error){
	data := []config.Data{
		config.Data{
			Metric:    "pgsql.maxconnection",
			Endpoint:  config.GetYaml().Ip,
			Tags:      "name=maxconnection",
			Value:     int(count),
			Timestamp: int(time.Now().Unix()),
			Step:      60,
		},
	}
	jsonStr, err := json.MarshalIndent(data, "", " ")
	if err!=nil{
		config.Error.Println(err)
	}
	config.Info.Println("formated: ", string(jsonStr))
	res,err:= util.Post(config.GetYaml().Api, jsonStr, "application/json")
	return res, err
}
