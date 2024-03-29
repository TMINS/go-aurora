package aurora

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

type Ctx struct {
	rw        *sync.RWMutex
	ar        *Aurora // Aurora 引用
	p         *proxy
	Response  http.ResponseWriter
	Request   *http.Request
	Args      map[string]interface{} //REST API 参数
	Attribute *sync.Map              //Context属性域
}

// Viper 获取项目配置实例,未启动配置则返回nil
func (c *Ctx) Viper() *viper.Viper {
	return c.ar.Viper()
}

// GetRoot 获取项目根路径
func (c *Ctx) GetRoot() string {
	return c.ar.projectRoot
}

// Mysql 获取注册的 默认 mysql
func (c *Ctx) Mysql() *gorm.DB {
	return get(c.ar.gorms, Mysql, 0)
}

// SQLite 获取注册的 默认 SQLite
func (c *Ctx) SQLite() *gorm.DB {
	return get(c.ar.gorms, SQLite, 0)
}

// PostgreSql 获取注册的 默认 PostgreSql
func (c *Ctx) PostgreSql() *gorm.DB {
	return get(c.ar.gorms, Postgresql, 0)
}

// SqlServer 获取注册的 默认 SqlServer
func (c *Ctx) SqlServer() *gorm.DB {
	return get(c.ar.gorms, SqlServer, 0)
}

// MysqlList 获取对应索引的 mysql
func (c *Ctx) MysqlList(index int) *gorm.DB {
	return get(c.ar.gorms, Mysql, index)
}

// SQLiteList 获取对应索引的 SQLite
func (c *Ctx) SQLiteList(index int) *gorm.DB {
	return get(c.ar.gorms, SQLite, index)
}

// PostgreSqlList 获取对应索引的 PostgreSql
func (c *Ctx) PostgreSqlList(index int) *gorm.DB {
	return get(c.ar.gorms, Postgresql, index)
}

// SqlServerList 获取对应索引的 SqlServer
func (c *Ctx) SqlServerList(index int) *gorm.DB {
	return get(c.ar.gorms, SqlServer, index)
}

// GoRedis 获取默认的 go-redis 客户端
func (c *Ctx) GoRedis() *redis.Client {
	return c.ar.goredis[0]
}

// GoRedisList 获取指定索引下的 go-redis 客户端
func (c *Ctx) GoRedisList(index int) *redis.Client {
	if index > len(c.ar.goredis) || index < 0 {
		//取客户端，超出索引边界给出提示
		c.ar.auroraLog.Warning("out of the storage index range of redis, the retrieval failed, please check the configuration information or parameters, a nil is returned")
		return nil
	}
	return c.ar.goredis[index]
}

// JSON 向浏览器输出json数据
func (c *Ctx) json(data interface{}) {
	s, b := data.(string) //返回值如果是json字符串或者直接是字符串，将不再转码,json 二次转码对原有的json格式会进行二次转义
	if b {
		c.Response.WriteHeader(http.StatusOK) // 写入200
		c.Response.Header().Set(contentType, c.ar.resourceMapType[".json"])
		_, err := c.Response.Write([]byte(s)) // 直接写入响应
		if err != nil {
			c.ar.auroraLog.Error(err.Error())
		}
		return
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		c.Response.WriteHeader(http.StatusBadGateway)
		_, err = c.Response.Write([]byte(err.Error()))
		c.ar.auroraLog.Error(err.Error())
		return
	}
	c.Response.Header().Set(contentType, c.ar.resourceMapType[".json"])
	c.Response.WriteHeader(http.StatusOK)
	_, err = c.Response.Write(marshal)
	if err != nil {
		c.ar.auroraLog.Error(err.Error())
	}
}

// RequestForward 内部路由转发，会携带 Ctx 本身进行转发，转发之后继续持有该 Ctx
func (c *Ctx) forward(path string) {
	c.ar.router.SearchPath(c.Request.Method, path, c.Response, c.Request, c)
}

// Redirect 发送重定向
func (c *Ctx) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}

// Message 主要用于在插件调用过程中出现 中断使用此方法对 用户进行响应消息
func (c *Ctx) Message(msg string) {
	c.Attribute.Store(plugin, msg)
}

func (c *Ctx) GetMessage(key interface{}) interface{} {
	load, ok := c.Attribute.Load(plugin)
	if !ok {
		return nil
	}
	return load
}
