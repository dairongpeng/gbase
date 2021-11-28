package gbase

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

// global config base on Viper
var v *viper.Viper
var gbaseOnce = sync.Once{}

func initViper() {
	v = viper.New()
	// 以CONFIG开头的环境变量会被加载进来。
	// 当使用 viper.Get(“apiversion”) 时，实际读取的环境变量是VIPER_APIVERSION。
	v.SetEnvPrefix("CONFIG")
	// ENV CONFIG_APP_NAME = "APP" can be accessed through Viper().GetString("app.name")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigFile("./config.yaml")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// Env config has high priority.
	// When file config and env config have the same key name, env key is used first
	v.AutomaticEnv()
	//viper.WatchConfig() 监控配置文件的变化。生产环境不建议打开
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	// 配置文件发生变更之后会调用的回调函数
	//	fmt.Println("Config file changed:", e.Name)
	//})
}

func Viper() *viper.Viper {
	if v == nil {
		gbaseOnce.Do(initViper)
	}
	return v
}
