package main

import (
	"context"
	"github.com/awensir/go-aurora/aurora"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {

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
