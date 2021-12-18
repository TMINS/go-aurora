package tests

import (
	"context"
	"fmt"
	"github.com/awensir/go-aurora/aurora"
	"github.com/awensir/go-aurora/aurora/pprofs"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"testing"
	"time"
)

type Mm1 struct {
}

func (m *Mm1) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm1")
	return true
}

func (m *Mm1) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm1")
}

func (m *Mm1) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm1")
}

type Mm2 struct {
}

func (m *Mm2) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm2")

	return true
}

func (m *Mm2) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm2")

}

func (m *Mm2) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm2")
}

type Mm3 struct {
}

func (m *Mm3) PreHandle(c *aurora.Ctx) bool {
	log.Println("PreHandle Mm3")
	return true
}

func (m *Mm3) PostHandle(c *aurora.Ctx) {
	log.Println("PostHandle Mm3")
}

func (m *Mm3) AfterCompletion(c *aurora.Ctx) {
	log.Println("AfterCompletion Mm3")
}

// 拦截器测试
func TestIntercept(t *testing.T) {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		log.Println("service..")
		return nil
	})

	a.RouteIntercept("/", &Mm1{}, &Mm2{}, &Mm3{})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}

// TestWebSocketClient 发起socket 测试客户端
func TestWebSocketClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080/", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "内部错误！")

	err = wsjson.Write(ctx, c, "Hello WebSocket Server")
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("接收到服务端响应：%v\n", v)

	c.Close(websocket.StatusNormalClosure, "")
}

// TestWebSocketServer websocket 服务端测试
func TestWebSocketServer(t *testing.T) {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		accept, err := websocket.Accept(c.Response, c.Request, nil)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
		defer cancel()
		var v interface{}
		err = wsjson.Read(ctx, accept, &v)
		if err != nil {
			return err
		}
		log.Printf("接收到客户端：%v\n", v)
		err = wsjson.Write(ctx, accept, "Hello WebSocket Client")
		if err != nil {
			return err
		}
		accept.Close(websocket.StatusNormalClosure, "")
		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}

// TestPprof 接口执行性能测试
func TestPprof(t *testing.T) {
	//获取 aurora 路由实例
	a := aurora.New()

	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {
		return nil
	})
	a.GET("/debug/pprof/heap", pprofs.Index)
	a.GET("/debug/pprof/cmdline", pprofs.Cmdline)
	a.GET("/debug/pprof/profile", pprofs.Profile)
	a.GET("/debug/pprof/symbol", pprofs.Symbol)
	a.GET("/debug/pprof/trace", pprofs.Trace)
	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()
}

// TestPlugins 插件中断测试
func TestPlugins(t *testing.T) {
	//获取 aurora 路由实例
	a := aurora.New()
	a.Plugin(func(ctx *aurora.Ctx) bool {
		fmt.Println("1")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("2")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("3")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("4")
		return true
	}, func(ctx *aurora.Ctx) bool {
		fmt.Println("5")
		ctx.Message("plugin error test!")
		return false
	})
	// GET 方法注册 web get请求
	a.GET("/", func(c *aurora.Ctx) interface{} {

		return nil
	})

	// 启动服务器 默认端口8080，更改端口号 a.Guide(”8081“) 即可
	a.Guide()

}
