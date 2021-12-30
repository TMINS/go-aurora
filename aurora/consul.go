package aurora

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"strconv"
	"strings"
)

const (
	address       = "address"       //HTTPAddrEnvName   			consul 服务地址
	token         = "token"         //HTTPTokenEnvName
	tokenFile     = "tokenFile"     //HTTPTokenFileEnvName
	ssl           = "ssl"           //HTTPSSLEnvName
	tlsServerName = "tlsServerName" //HTTPTLSServerName
	cafile        = "cafile"        //HTTPCAFile
	capath        = "capath"        //HTTPCAPath
	clientCert    = "clientCert"    //HTTPClientCert
	clientKey     = "clientKey"     //HTTPClientKey
	namespace     = "namespace"     //HTTPNamespaceEnvName
	verify        = "verify"        //HTTPSSLVerifyEnvName
	auth          = "auth"          //HTTPAuthEnvName
)

/*
	consul 模块
*/

//  consulConfig 读取配置文件加载
func (a *Aurora) consulConfig() {
	if a.cnf == nil {
		return
	}
	var (
		b   bool
		err error
	)
	c := a.cnf.Get("aurora.consul")
	if c == nil {
		return
	}
	option, b := c.(map[string]interface{})
	if !b {
		return
	}
	config := api.DefaultConfig()
	if v, b := option["address"]; b {
		config.Address = v.(string)
	}
	if v, b := option["tlsServerName"]; b {
		config.TLSConfig.Address = v.(string)
	}
	if v, b := option["cafile"]; b {
		config.TLSConfig.CertFile = v.(string)
	}
	if v, b := option["capath"]; b {
		config.TLSConfig.CAPath = v.(string)
	}
	if v, b := option["clientCert"]; b {
		config.TLSConfig.CertFile = v.(string)
	}
	if v, b := option["clientKey"]; b {
		config.TLSConfig.KeyFile = v.(string)
	}
	if v, b := option["namespace"]; b {
		config.Namespace = v.(string)
	}
	if v, b := option["tokenFile"]; b {
		config.TokenFile = v.(string)
	}
	if v, b := option["token"]; b {
		config.Token = v.(string)
	}

	if v, b := option["auth"]; b {
		var username, password string
		if strings.Contains(v.(string), ":") {
			split := strings.SplitN(v.(string), ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = v.(string)
		}
		config.HttpAuth = &api.HttpBasicAuth{
			Username: username,
			Password: password,
		}
	}

	if v, b := option["ssl"]; b {
		enabled, err := strconv.ParseBool(v.(string))
		if err != nil {
			fmt.Errorf("could not parse ssl %s error", err)
		}
		if enabled {
			config.Scheme = "https"
		}
	}
	if v, b := option["verify"]; b {
		doVerify, err := strconv.ParseBool(v.(string))
		if err != nil {
			fmt.Errorf("could not parse verify error: %s ", err)
		}
		if !doVerify {
			config.TLSConfig.InsecureSkipVerify = true
		}
	}
	a.consulClient, err = api.NewClient(config)
	if err != nil {
		panic(err)
	}
	a.Start()
}

func (a *Aurora) Start() {
	//拿到代理实例
	agent := a.consulClient.Agent()
	check := new(api.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("%s://%s:%s/actuator/health", "http", "localhost", a.port)
	a.GET("/actuator/health", func(c *Ctx) interface{} {
		return "ok"
	})
	check.Timeout = "5s"
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "20s" // 故障检查失败30s后 consul自动将注册服务删除
	err := agent.ServiceRegister(&api.AgentServiceRegistration{
		ID:      "IDS",
		Name:    a.name,
		Address: ServerIp(),
		Port:    a.Port(),
		Check:   check,
	})

	if err != nil {
		fmt.Println("注册失败,error:", err.Error())
		return
	}
	fmt.Println("注册成功")
}

// ServerIp 获取服务器ip地址信息
func ServerIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
func (a *Aurora) Port() int {
	atoi, err := strconv.Atoi(a.port)
	if err != nil {
		return 0
	}
	return atoi
}
