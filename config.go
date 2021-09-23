package gbase

import (
	"github.com/spf13/viper"
	"strings"
)

// global config base on Viper
var v *viper.Viper

func init() {
	v = viper.New()
	v.SetEnvPrefix("CONFIG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

func Viper() *viper.Viper {
	return v
}
