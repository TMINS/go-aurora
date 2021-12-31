package aurora

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path"
)

const FILE = "application.yml"

/*
	viper 配置文件实例
	提供Aurora 默认的配置实例
	默认读取配置文件的位置为根目录 application.yml
	默认配置文件项，优先于api配置项
*/

type auroraConfig struct {
	config interface{}
}

// viperConfig 配置并加载 application.yml 配置文件
func (a *Aurora) viperConfig(p ...string) {
	if a.cnf != nil {
		return
	}
	a.cnf = viper.New() //创建配置文件实例
	cnf := make([]string, 0)
	cnf = append(cnf, a.projectRoot) //添加项目根路径
	if p != nil {
		cnf = append(cnf, p...)
	} else {
		//检查默认配置文件 是否存在 不存在则不做任何加载
		pat := path.Join(a.projectRoot, FILE)
		_, err := os.Lstat(pat)
		if err != nil {
			if os.IsNotExist(err) {
				a.message <- fmt.Sprintf("配置文件不存在，加载失败")
			}
			return
		}
	}
	cnf = append(cnf, FILE)
	conf := path.Join(cnf...)
	a.cnf.SetConfigFile(conf)
	err := a.cnf.ReadInConfig()
	if err != nil {
		a.message <- err.Error()
		return
	}
	a.message <- fmt.Sprint("配置文件加载成功.")
}

// Viper 获取 Aurora viper实例
func (a *Aurora) Viper() *viper.Viper {
	return a.cnf
}
