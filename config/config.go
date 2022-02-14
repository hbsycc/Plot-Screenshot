package config

import (
	"a.resources.cc/model"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

var config = new(model.Config)

// SetConfig
// @Description: 读取配置文件并初始化配置参数
// @return err
func SetConfig() (err error) {
	configFile, err := os.Open("./config.json")
	if err != nil {
		return
	} else {
		defer func(configFile *os.File) {
			_ = configFile.Close()
		}(configFile)
	}

	c, err := ioutil.ReadAll(configFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(c, &config)
	for i, s := range config.Media.Ext {
		config.Media.Ext[i] = strings.ToLower(s)
	}
	if err != nil {
		return
	}

	return
}

// GetConfig
// @Description: 获取配置参数
func GetConfig() model.Config {
	return *config
}
