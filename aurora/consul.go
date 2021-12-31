package aurora

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"strconv"
	"strings"
)

const (
	address       = "address" //HTTPAddrEnvName   			consul 服务地址
	scheme        = "scheme"
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

	serviceType = "type"       //服务类型
	timeout     = "timeout"    //健康检查超时时间
	interval    = "interval"   //健康检查间隔时间
	deregister  = "deregister" //最大超时时间 注销服务

	defaultTimeout    = "5s"  //默认超时时间
	defaultInterval   = "4s"  //默认检查间隔时间
	defaultDeregister = "30s" //默认最大超时注销时间

	HTTP = iota
	gRPC
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
		b bool
	)
	c := a.cnf.Get("aurora.consul") //读取consul配置信息
	if c == nil {
		return
	}
	option, b := c.(map[string]interface{})
	if !b {
		return
	}
	if addrs, v := option[address]; v {
		config := DefaultConsulConfig(option)
		switch addrs.(type) {
		case string:
			client := ConsulApiClientConfig(addrs.(string), config)
			//初始化api客户端完毕，开始向consul服务器注册服务
			a.ConsulServiceRegister(client, config, option)
			a.consuls = append(a.consuls, client)
		case []string:
			for _, v := range addrs.([]string) {
				client := ConsulApiClientConfig(v, config)
				//初始化api客户端完毕，开始向consul服务器注册服务
				a.ConsulServiceRegister(client, config, option)
				a.consuls = append(a.consuls, client)
			}
		}
	}
}

// ConsulApiClientConfig 配置并返回一个 consul api实例，多个地址公用一套配置项
func ConsulApiClientConfig(address string, config *api.Config) *api.Client {
	config.Address = address //配置默认提供的config 的address 生成一个client 实例
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	return client
}

func (a *Aurora) ConsulServiceRegister(client *api.Client, config *api.Config, option map[string]interface{}) {
	//创建健康检查
	check := &api.AgentServiceCheck{}
	if v, b := option[timeout]; b {
		check.Timeout = v.(string)
	} else {
		check.Timeout = defaultTimeout
	}

	if v, b := option[interval]; b {
		check.Interval = v.(string)
	} else {
		check.Interval = defaultInterval
	}

	if v, b := option[deregister]; b {
		check.DeregisterCriticalServiceAfter = v.(string)
	} else {
		check.DeregisterCriticalServiceAfter = defaultDeregister
	}
	//配置服务健康检查接口
	if v, b := option[serviceType]; b {
		switch v.(string) {
		case "HTTP", "http":
			check.HTTP = fmt.Sprintf("%s://%s:%d/consul/health", config.Scheme, ServerIp(), a.port) //配置健康检查接口
			a.GET("/consul/health", HTTPCheck)
		case "GRPC", "grpc", "gRPC":
			//先检查grpc 是否进行配置
			if a.grpc == nil {
				panic("服务没有配置grpc")
			}
			check.GRPC = fmt.Sprintf("%s:%d", ServerIp(), a.port) //配置健康检查接口
			check.GRPCUseTLS = false
			grpc_health_v1.RegisterHealthServer(a.grpc, health.NewServer())
		}
	} else {
		//默认采用http服务
		check.HTTP = fmt.Sprintf("%s://%s:%d/consul/agent/health", config.Scheme, ServerIp(), a.Port()) //配置健康检查接口
		a.GET("/consul/agent/health", HTTPCheck)
	}
	//创建准备注册的服务信息，服务的基本信息基于本程序的本机信息
	service := &api.AgentServiceRegistration{}

	service.Name = a.name        //Name 属性标识在consul中的服务名称
	service.ID = a.name          //ID属性是基于Name属性下面的编号,使用的时候不应该出现重复
	service.Check = check        //Check 属性用于配置，相对的Checks可以配置多个
	service.Address = ServerIp() //设置服务地址信息
	service.Port = a.Port()      //设置端口信息

	agent := client.Agent()

	err := agent.ServiceRegister(service)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// ConsulHttpServiceRegister 注册HTTP类型的服务
func ConsulHttpServiceRegister() {

}

// ConsulGrpcServiceRegister 注册gRPC类型的服务
func ConsulGrpcServiceRegister() {
	newServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(newServer, health.NewServer())
}

// DefaultConsulConfig 配置一个默认的Consul配置项，该配置实例address为默认值(localhost)，即使是集群地址也是每个地址采用一样的配置，该config 主要用来复用注册集群
func DefaultConsulConfig(option map[string]interface{}) *api.Config {
	config := api.DefaultConfig()
	if v, b := option[scheme]; b {
		config.Scheme = v.(string)
	}
	if v, b := option[tlsServerName]; b {
		config.TLSConfig.Address = v.(string)
	}
	if v, b := option[cafile]; b {
		config.TLSConfig.CertFile = v.(string)
	}
	if v, b := option[capath]; b {
		config.TLSConfig.CAPath = v.(string)
	}
	if v, b := option[clientCert]; b {
		config.TLSConfig.CertFile = v.(string)
	}
	if v, b := option[clientKey]; b {
		config.TLSConfig.KeyFile = v.(string)
	}
	if v, b := option[namespace]; b {
		config.Namespace = v.(string)
	}
	if v, b := option[tokenFile]; b {
		config.TokenFile = v.(string)
	}
	if v, b := option[token]; b {
		config.Token = v.(string)
	}
	if v, b := option[auth]; b {
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
	if v, b := option[ssl]; b {
		enabled, err := strconv.ParseBool(v.(string))
		if err != nil {
			fmt.Errorf("could not parse ssl %s error", err)
		}
		if enabled {
			config.Scheme = "https"
		}
	}
	if v, b := option[verify]; b {
		doVerify, err := strconv.ParseBool(v.(string))
		if err != nil {
			fmt.Errorf("could not parse verify error: %s ", err)
		}
		if !doVerify {
			config.TLSConfig.InsecureSkipVerify = true
		}
	}
	return config
}

// HTTPCheck HTTP健康检查
func HTTPCheck(c *Ctx) interface{} {
	return "ok"
}
