package aurora

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/awensir/aurora-email/email"
	"github.com/awensir/go-aurora/aurora/req"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

//2022/01/19 重新构建 处理函数

// HttpRequest 作为处理函数的主要参数,其中会初始化一个上下文参数，如果无法初始化上下文参数则不会进入到处理方法,为了不重复的构建api HttpRequest将对Ctx进行封装
// HttpRequest 解析的参数默认类型：
// 数字类型: float64
// json数据：map[string]interface{}
// HttpRequest 只用于封装ctx的对外调用和数据解析的承载
type HttpRequest map[string]interface{}

type HttpHandle interface {
	Hadnle(HttpRequest) interface{}
}

type Handel func(HttpRequest) interface{}

func (handel Handel) Hadnle(hq HttpRequest) interface{} {
	defer func(hq HttpRequest) {
		ctx := hq[req.Ctx].(*Ctx)
		v := recover()
		if v != nil {
			// Serve 处理器发生 panic 等严重错误处理，给调用者返回 500 并返回错误描述
			switch v.(type) {
			case string:
				http.Error(ctx.Response, v.(string), 500)
			case error:
				http.Error(ctx.Response, v.(error).Error(), 500)
			default:
				marshal, err := json.Marshal(v)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				http.Error(ctx.Response, string(marshal), 500)
			}
			return
		}
	}(hq)
	return handel(hq)
}

//封装基础组件调用

func (r HttpRequest) Mysql() *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.Mysql()
}

func (r HttpRequest) PostgreSql() *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.PostgreSql()
}

func (r HttpRequest) SQLite() *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.SQLite()
}

func (r HttpRequest) SqlServer() *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.SqlServer()
}

func (r HttpRequest) MysqlList(index int) *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.MysqlList(index)
}

func (r HttpRequest) PostgreSqlList(index int) *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.PostgreSqlList(index)
}

func (r HttpRequest) SQLiteList(index int) *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.SQLiteList(index)
}

func (r HttpRequest) SqlServerList(index int) *gorm.DB {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.SqlServerList(index)
}

//保存文件
func (r HttpRequest) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.SaveUploadedFile(file, dst)
}

//获取邮件客户端
func (r HttpRequest) Email() *email.Client {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.Email()
}

//获取redis 客户端，redis 待升级支持多个
func (r HttpRequest) GoRedis() *redis.Client {
	ctx := r[req.Ctx].(*Ctx)
	return ctx.GoRedis()
}
