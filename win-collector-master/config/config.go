package config

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/n9e/win-collector/stra"
	"github.com/n9e/win-collector/sys"
	"github.com/n9e/win-collector/sys/identity"

	"github.com/didi/nightingale/src/toolkits/logger"
	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

type ConfYaml struct {
	Identity identity.IdentitySection `yaml:"identity"`
	Logger   logger.LoggerSection     `yaml:"logger"`
	Stra     stra.StraSection         `yaml:"stra"`
	Sys      sys.SysSection           `yaml:"sys"`
}

var (
	Config   *ConfYaml
	lock     = new(sync.RWMutex)
	Endpoint string
	Cwd      string
)

// Get configuration file
func Get() *ConfYaml {
	lock.RLock()
	defer lock.RUnlock()
	return Config
}

func Parse(conf string) error {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	lock.Lock()
	defer lock.Unlock()

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	viper.SetDefault("worker", map[string]interface{}{
		"workerNum":    10,
		"queueSize":    1024000,
		"pushInterval": 5,
		"waitPush":     0,
	})

	viper.SetDefault("stra", map[string]interface{}{
		"enable":   true,
		"timeout":  1000,
		"interval": 10, //采集策略更新时间
		"portPath": "/home/n9e/etc/port",
		"procPath": "/home/n9e/etc/proc",
		"logPath":  "/home/n9e/etc/log",
		"api":      "/api/portal/collects/",
	})

	viper.SetDefault("sys", map[string]interface{}{
		"timeout":  1000, //请求超时时间
		"interval": 10,   //基础指标上报周期
		"plugin":   "/home/n9e/plugin",
	})

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("Unmarshal %v", err)
	}

	return nil
}