package aurora

import "github.com/awensir/aurora-email/email"

func (a *Aurora) loadEmail() {
	if a.cnf == nil {
		//如果配置文件没有加载成功，将不做任何事情
		return
	}
	EmailCinfig := a.cnf.GetStringMapString("aurora.email")
	var user, passwd, host string
	if v, b := EmailCinfig["user"]; b {
		if v != "" {
			user = v
		}
	}
	if v, b := EmailCinfig["password"]; b {
		if v != "" {
			passwd = v
		}
	}
	if v, b := EmailCinfig["host"]; b {
		if v != "" {
			host = v
		}
	}
	if user == "" || passwd == "" || host == "" {
		return
	}
	a.email = email.NewClient(user, passwd, host) //初始化email
}

// Email 获取email 客户端，若无相关的配置信息，则返回nil
func (a *Aurora) Email() *email.Client {
	return a.email
}

func (c *Ctx) Email() *email.Client {
	return c.ar.email
}
