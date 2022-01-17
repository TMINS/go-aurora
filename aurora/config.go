package aurora

import (
	"bytes"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"os"
	"path"
	"sync"
	"time"
)

const FILE = "application.yml"

/*
	viper 配置文件实例
	提供Aurora 默认的配置实例
	默认读取配置文件的位置为根目录 application.yml
	默认配置文件项，优先于api配置项
*/

type ConfigCenter struct {
	cnf *viper.Viper
	rw  *sync.RWMutex
}

func (c *ConfigCenter) GetStringMapString(key string) map[string]string {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnf.GetStringMapString(key)
}

func (c *ConfigCenter) Get(key string) interface{} {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnf.Get(key)
}

func (c *ConfigCenter) GetStringMap(key string) map[string]interface{} {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnf.GetStringMap(key)
}

func (c *ConfigCenter) GetString(key string) string {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnf.GetString(key)
}

func (c *ConfigCenter) GetStringMapStringSlice(key string) map[string][]string {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnf.GetStringMapStringSlice(key)
}

// viperConfig 配置并加载 application.yml 配置文件
func (a *Aurora) viperConfig(p ...string) {
	if a.config == nil {
		a.config = &ConfigCenter{
			cnf: viper.New(),
			rw:  &sync.RWMutex{},
		}
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
				//没有加载到配置文件 给出警告
				a.auroraLog.Warning(fmt.Sprintf("no configuration file information was loaded because the default application.yml was not found."))
			}
			return
		}
	}
	cnf = append(cnf, FILE)
	conf := path.Join(cnf...)
	a.config.cnf.SetConfigFile(conf)
	err := a.config.cnf.ReadInConfig()
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}
	//开始检查加载远程配置中心
	NacosConfig := a.config.cnf.GetStringMap("remote.config.nacos")
	a.remoteConfigs(NacosConfig)

	a.auroraLog.Info("the configuration file is loaded successfully.")
}

// Viper 获取 Aurora viper实例
func (a *Aurora) Viper() *viper.Viper {
	return a.cnf
}

// 远程配置中心读取 在开发中，未测试
func (a *Aurora) remoteConfigs(remote map[string]interface{}) {
	if remote == nil {
		//没有读取到远程配置中心
		return
	}
	var ServerConfig []constant.ServerConfig
	var ClientConfig *constant.ClientConfig
	dataid := remote["DataId"].(string)
	gropu := remote["Group"].(string)
	//配置服务器
	if v, b := remote["server"]; b {
		servers, f := v.([]interface{})
		if f {
			//如果失败则，配置参数是存在问题的，nacos的服务器参数配置传递参数是数组方式，并且至少一个服务器实例
			a.auroraLog.Fatal("the server parameter information in the nacos remote configuration center is incorrectly configured. check whether the configuration is in the array mode.")
		}
		ServerConfig = make([]constant.ServerConfig, 0)
		for _, s := range servers {
			args := s.(map[string]interface{})
			server := constant.ServerConfig{}
			if field, t := args["IpAddr"]; t {
				server.IpAddr = field.(string)
			}
			if field, t := args["Prot"]; t {
				server.Port = field.(uint64)
			}
			if field, t := args["Scheme"]; t {
				server.Scheme = field.(string)
			}
			if field, t := args["ContextPath"]; t {
				server.ContextPath = field.(string)
			}
			ServerConfig = append(ServerConfig, server)
		}
	}

	//初始化 客户端配置
	if v, b := remote["client"]; b {
		ClientConfig = &constant.ClientConfig{}
		args := v.(map[string]interface{})
		if field, t := args["NamespaceId"]; t {
			ClientConfig.NamespaceId = field.(string)
		}
		if field, t := args["TimeoutMs"]; t {
			ClientConfig.TimeoutMs = field.(uint64)
		}
		if field, t := args["ListenInterval"]; t {
			ClientConfig.ListenInterval = field.(uint64)
		}
		if field, t := args["BeatInterval"]; t {
			ClientConfig.BeatInterval = field.(int64)
		}
		if field, t := args["AppName"]; t {
			ClientConfig.AppName = field.(string)
		}
		if field, t := args["Endpoint"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["RegionId"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["AccessKey"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["SecretKey"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["OpenKMS"]; t {
			ClientConfig.OpenKMS = field.(bool)
		}
		if field, t := args["CacheDir"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["UpdateThreadNum"]; t {
			ClientConfig.UpdateThreadNum = field.(int)
		}
		if field, t := args["NotLoadCacheAtStart"]; t {
			ClientConfig.NotLoadCacheAtStart = field.(bool)
		}
		if field, t := args["UpdateCacheWhenEmpty"]; t {
			ClientConfig.UpdateCacheWhenEmpty = field.(bool)
		}
		if field, t := args["Username"]; t {
			ClientConfig.Username = field.(string)
		}
		if field, t := args["Password"]; t {
			ClientConfig.Password = field.(string)
		}
		if field, t := args["LogDir"]; t {
			ClientConfig.LogDir = field.(string)
		}
		if field, t := args["RotateTime"]; t {
			ClientConfig.RotateTime = field.(string)
		}
		if field, t := args["MaxAge"]; t {
			ClientConfig.MaxAge = field.(int64)
		}
		if field, t := args["LogLevel"]; t {
			ClientConfig.LogLevel = field.(string)
		}
		if field, t := args["LogSampling"]; t {
			arg := field.(map[string]interface{})
			ClientConfig.LogSampling = &constant.ClientLogSamplingConfig{}
			if f, is := arg["Initial"]; is {
				ClientConfig.LogSampling.Initial = f.(int)
			}
			if f, is := arg["Thereafter"]; is {
				ClientConfig.LogSampling.Thereafter = f.(int)
			}
			if f, is := arg["Tick"]; is {
				ClientConfig.LogSampling.Tick = f.(time.Duration)
			}
		}
		if field, t := args["ContextPath"]; t {
			ClientConfig.ContextPath = field.(string)
		}
	}
	if ServerConfig == nil || ClientConfig == nil {
		return
	}
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  ClientConfig,
		ServerConfigs: ServerConfig,
	})
	if err != nil {
		a.auroraLog.Fatal("nacos remote configuration failed to load error:", err.Error())
		return
	}

	//开始读取配置文件
	config, err := client.GetConfig(vo.ConfigParam{
		DataId: dataid,
		Group:  gropu,
	})
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}

	//初次加载远程配置
	buf := bytes.NewBufferString(config)
	a.cnf.SetConfigType("yml")
	err = a.cnf.ReadConfig(buf)
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}

	//启动远程配置监听
	err = client.ListenConfig(vo.ConfigParam{
		DataId:   dataid,
		Group:    gropu,
		OnChange: a.refreshConfig,
	})
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}
}

// 重新加载远程配置文件
func (a *Aurora) refreshConfig(namespace, group, dataId, data string) {
	//刷新服务运行配置
	buf := bytes.NewBufferString(data)
	err := a.config.cnf.ReadConfig(buf) //需要分装加锁
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}
	a.auroraLog.Info("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
}
