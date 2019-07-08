package conf

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sync"

	utils "trade_api/src/main/util"
	"truxing/commons/conf"
)

var lock sync.Mutex
var baseConfig *Config

type Config struct {
	RootDir  string
	Resource *conf.ResourceConfig
}

var (
	sep = string(os.PathSeparator)
)

var (
	port string
	mode string
)

func GetPort() string {
	return port
}

func GetMode() string {
	return mode
}

func init() {
	flag.StringVar(&port, "port", ":25433", "启用的端口号，默认为25433")
	flag.StringVar(&mode, "mode", "dev", "配置文件模式")
	flag.Parse()
}

// InitConfig must run before the server start
func InitConfig() {
	lock.Lock()
	defer lock.Unlock()
	curDir := CurrentDir()
	var f string
	if mode == "dev" {
		f = "json" + sep + "base.json"
	} else {
		f = "json_" + mode + sep + "base.json"
	}
	f = filepath.Join(curDir, "conf", f)
	s, err := utils.ExtendFile(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(s), &baseConfig)
	if err != nil {
		panic(err)
	}
	conf.InitResConfig(baseConfig.Resource)
	baseConfig.RootDir = filepath.Join(curDir, ".."+sep)
}

func LoadResource() *conf.ResourceConfig {
	curDir := CurrentDir()
	var f string
	if mode == "dev" {
		f = "json" + sep + "resource.json"
	} else {
		f = "json_" + mode + sep + "resource.json"
	}
	f = filepath.Join(curDir, "conf", f)
	s, err := utils.ExtendFile(f)
	if err != nil {
		panic(err)
	}

	resource := new(conf.ResourceConfig)
	err = json.Unmarshal([]byte(s), resource)
	if err != nil {
		panic(err)
	}
	return resource
}

//获取调用者的路径
func CurrentDir() string {
	dir, _ := os.Getwd()
	return dir
}

func Conf() Config {
	if baseConfig == nil {
		InitConfig()
	}
	return *baseConfig
}
