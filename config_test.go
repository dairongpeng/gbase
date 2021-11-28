package gbase

import (
	"fmt"
	"os"
	"testing"
)

func Test_Config(t *testing.T) {
	// 使用环境变量
	os.Setenv("CONFIG_USER_SECRET_ID", "QLdywI2MrmDVjSSv6e95weNRvmteRjfKAuNV")
	os.Setenv("CONFIG_USER_SECRET_KEY", "bVix2WBv0VPfrDrvlLWrhEdzjLpPCNYb")

	//viper.AutomaticEnv()                                             // 读取环境变量
	//viper.SetEnvPrefix("VIPER")                                      // 设置环境变量前缀：VIPER_，如果是viper，将自动转变为大写。
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // 将viper.Get(key) key字符串中'.'和'-'替换为'_'
	//viper.BindEnv("user.secret-key")
	//viper.BindEnv("user.secret-id", "USER_SECRET_ID") // 绑定环境变量名到key
	initViper()
	fmt.Println(Viper().GetString("USER.SECRET.ID"))
}
