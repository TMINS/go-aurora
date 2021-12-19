package aurora

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
)

const CNF_FILE = "application.yml"

/*
	viper 配置文件实例
	提供Aurora 默认的配置实例
	默认读取配置文件的位置为根目录 application.yml
*/

// Opt 配置选项参数
type Opt func() map[string]interface{}

// Configuration 第三方配置加载，传入指定配置项，返回被配置的对象
type Configuration func(opt Opt) interface{}

// Config 存储一个 配置实例，以便多次配置
type Config struct {
	Opt //配置参数
	Configuration
}

// store 加载并返回配置实例
func (o *Config) store() interface{} {
	return o.Configuration(o.Opt)
}

// ViperConfig 配置并加载 application.yml 配置文件
func (a *Aurora) ViperConfig(p ...string) {
	a.cnf = viper.New()
	cnf := make([]string, 0)
	cnf = append(cnf, a.projectRoot)
	//cnfpath:=path.Join(a.projectRoot,CNF_FILE)
	if p != nil {
		cnf = append(cnf, p...)
	} else {
		//检查默认配置文件 是否存在 不存在则不做任何加载
		pat := path.Join(a.projectRoot, CNF_FILE)
		_, err := os.Lstat(pat)
		if err != nil {
			if os.IsNotExist(err) {
				a.message <- fmt.Sprintf("No configuration file loaded")
			}
			return
		}
	}
	cnf = append(cnf, CNF_FILE)
	conf := path.Join(cnf...)
	a.cnf.SetConfigType("yml")
	a.cnf.SetConfigFile(conf)
	err := a.cnf.ReadInConfig()
	if err != nil {
		a.message <- err.Error()
		return
	}
	a.message <- fmt.Sprint("Aurora Configuration file loaded successfully")
}

// Viper 获取 Aurora viper实例
func (a *Aurora) Viper() *viper.Viper {
	return a.cnf
}
