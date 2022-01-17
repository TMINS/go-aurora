package aurora

import (
	"bytes"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"io"
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

// ConfigCenter 配置中心 的读写锁主要用来解决分布式配置的动态刷新配置，和以后存在的并发读取配置和修改
// 对于修改配置数据库连接信息或者需要重新初始化的配置项这些无法起到同步更新的效果只能保持配置信息是最新的（需要重新初始化的配置建议重启服务）
// 对被配置的使用实例没有并发安全的效果
type ConfigCenter struct {
	cnf *viper.Viper
	rw  *sync.RWMutex
}

func (c *ConfigCenter) SetConfigFile(in string) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	c.cnf.SetConfigFile(in)
}

func (c *ConfigCenter) SetConfigType(in string) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	c.cnf.SetConfigType(in)
}

func (c *ConfigCenter) ReadInConfig() error {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnf.ReadInConfig()
}

func (c *ConfigCenter) ReadConfig(in io.Reader) error {
	//写锁
	c.rw.Lock()
	defer c.rw.Unlock()
	return c.cnf.ReadConfig(in)
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

func (c *ConfigCenter) GetStringSlice(key string) []string {
	c.rw.RLock()
	defer c.rw.RUnlock()
	return c.cnf.GetStringSlice(key)
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
	a.config.SetConfigFile(conf)
	err := a.config.ReadInConfig()
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}
	//开始检查加载远程配置中心
	NacosConfig := a.config.GetStringMap("aurora.nacos")
	a.remoteConfigs(NacosConfig)

	a.auroraLog.Info("the configuration file is loaded successfully.")
}

// Viper 获取 Aurora viper实例
func (a *Aurora) Viper() *viper.Viper {
	return a.cnf
}

// 远程配置中心读取 在开发中，测试能够读取到配置
func (a *Aurora) remoteConfigs(remote map[string]interface{}) {
	if len(remote) == 0 {
		//没有读取到远程配置中心
		return
	}
	var ServerConfig []constant.ServerConfig
	var ClientConfig *constant.ClientConfig
	dataid := remote["dataid"].(string)
	gropu := remote["group"].(string)
	if dataid == "" || gropu == "" {
		a.auroraLog.Fatal("please check whether the dataid or group configuration in the configuration file is configured correctly")
	}
	//配置服务器
	if v, b := remote["serverconfig"]; b {
		servers, f := v.([]interface{})
		if !f {
			//如果失败则，配置参数是存在问题的，nacos的服务器参数配置传递参数是数组方式，并且至少一个服务器实例
			a.auroraLog.Fatal("the server parameter information in the nacos remote configuration center is incorrectly configured. check whether the configuration is in the array mode.")
		}
		ServerConfig = make([]constant.ServerConfig, 0)
		for _, s := range servers {
			args := s.(map[interface{}]interface{})
			server := constant.ServerConfig{}
			if field, t := args["IpAddr"]; t {
				server.IpAddr = field.(string)
			}
			if field, t := args["Prot"]; t {
				server.Port = uint64(field.(int))
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
	if v, b := remote["clientconfig"]; b {
		ClientConfig = &constant.ClientConfig{}
		args := v.(map[string]interface{})
		if field, t := args["namespaceid"]; t {
			ClientConfig.NamespaceId = field.(string)
		}
		if field, t := args["timeoutms"]; t {
			ClientConfig.TimeoutMs = uint64(field.(int))
		}
		if field, t := args["listeninterval"]; t {
			ClientConfig.ListenInterval = uint64(field.(int))
		}
		if field, t := args["beatinterval"]; t {
			ClientConfig.BeatInterval = field.(int64)
		}
		if field, t := args["appname"]; t {
			ClientConfig.AppName = field.(string)
		}
		if field, t := args["endpoint"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["regionid"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["accesskey"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["secretkey"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["openkms"]; t {
			ClientConfig.OpenKMS = field.(bool)
		}
		if field, t := args["cachedir"]; t {
			ClientConfig.Endpoint = field.(string)
		}
		if field, t := args["updatethreadnum"]; t {
			ClientConfig.UpdateThreadNum = field.(int)
		}
		if field, t := args["notLoadcacheatstart"]; t {
			ClientConfig.NotLoadCacheAtStart = field.(bool)
		}
		if field, t := args["updatecachewhenempty"]; t {
			ClientConfig.UpdateCacheWhenEmpty = field.(bool)
		}
		if field, t := args["username"]; t {
			ClientConfig.Username = field.(string)
		}
		if field, t := args["password"]; t {
			ClientConfig.Password = field.(string)
		}
		if field, t := args["logdir"]; t {
			ClientConfig.LogDir = field.(string)
		}
		if field, t := args["rotatetime"]; t {
			ClientConfig.RotateTime = field.(string)
		}
		if field, t := args["maxage"]; t {
			ClientConfig.MaxAge = int64(field.(int))
		}
		if field, t := args["loglevel"]; t {
			ClientConfig.LogLevel = field.(string)
		}
		if field, t := args["logsampling"]; t {
			arg := field.(map[string]interface{})
			ClientConfig.LogSampling = &constant.ClientLogSamplingConfig{}
			if f, is := arg["initial"]; is {
				ClientConfig.LogSampling.Initial = f.(int)
			}
			if f, is := arg["thereafter"]; is {
				ClientConfig.LogSampling.Thereafter = f.(int)
			}
			if f, is := arg["tick"]; is {
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
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  ClientConfig,
			ServerConfigs: ServerConfig,
		},
	)
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
	a.auroraLog.Info("remote configuration loaded successfully..")
	//初次加载远程配置
	buf := bytes.NewBufferString(config)
	a.cnf.SetConfigType("yml")
	err = a.cnf.ReadConfig(buf)
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}
	a.auroraLog.Info("start remote configuration file modification monitoring")
	//启动远程配置监听
	err = client.ListenConfig(vo.ConfigParam{
		DataId:   dataid,
		Group:    gropu,
		OnChange: a.refreshNacosConfig,
	})
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}
}

// 重新加载远程配置文件
func (a *Aurora) refreshNacosConfig(namespace, group, dataId, data string) {
	//刷新服务运行配置
	a.auroraLog.Info("refresh configuration...")
	buf := bytes.NewBufferString(data)
	err := a.config.ReadConfig(buf)
	if err != nil {
		a.auroraLog.Fatal(err.Error())
		return
	}
	a.auroraLog.Info("namespace:", namespace, ", config changed group:"+group+", dataId:"+dataId+", refresh successfully")
}
