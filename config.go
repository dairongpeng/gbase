package base

import (
	"github.com/spf13/viper"
	"strings"
)

// global config base on Viper
var config *viper.Viper

func init() {
	config = viper.New()
	config.SetEnvPrefix("BASE")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
}

func V() *viper.Viper {
	return config
}
