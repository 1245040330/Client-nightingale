package fengjiePlugins

import (
	"fmt"
	"github.com/didi/nightingale/src/modules/collector/config"
	"github.com/didi/nightingale/src/modules/collector/nginxInfo"
	"github.com/didi/nightingale/src/modules/collector/pgsqlInfo"
	"github.com/didi/nightingale/src/modules/collector/redisInfo"
	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/runner"
	"os"
	"path"
	"time"
)



func Detect() {
	//读取静态资源
	conf:=aconf()

	pconf(conf)
	detect()
	go loopDetect()
}

func loopDetect() {
	for {
		time.Sleep(time.Second * 60)
		detect()
	}
}

func detect() {
	conf:=config.GetYaml()
	if(conf.Pgsql.Dbname!=""){
		//pg数据获取
		resPg,err:=pgsqlInfo.SqlInfo()
		if err!=nil{
			config.Error.Println(err)
		}
		config.Info.Println("pg数据提交返回\n"+resPg)
	}
	if(conf.Redis.Addr!=""){
		//redis数据获取
		resRedis,err:=redisInfo.RedisInfo()
		if err!=nil{
			config.Error.Println(err)
		}
		config.Info.Println("redis数据提交返回\n"+resRedis)
	}
	if(conf.NginxStatus!=""){
		//nginx 获取数据
		resNginx,err:=nginxInfo.NginxInfo()
		if err!=nil{
			config.Error.Println(err)
		}
		config.Info.Println("nginx数据提交返回\n"+resNginx)
	}


}

func aconf() string{

	conf := path.Join(runner.Cwd, "etc", "configuration.local.yml")
	if file.IsExist(conf) {
		return conf
	}

	conf= path.Join(runner.Cwd, "etc", "configuration.yml")
	if file.IsExist(conf) {
		return conf
	}
	fmt.Println("no configuration file for sender")
	os.Exit(1)
	return ""
}
func pconf(conf1 string) {
	if err := config.ParseConfig(conf1); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	} else {
		fmt.Println("parse configuration file:", conf1)
	}
}