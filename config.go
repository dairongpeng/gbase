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

	initLog()
}

func Viper() *viper.Viper {
	if v == nil {
		gbaseOnce.Do(initViper)
	}
	return v
}
