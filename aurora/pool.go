package aurora

import (
	"github.com/awensir/go-aurora/aurora/option"
	"sync"
)

/*
	通过pool的方式添加 配置，能够保证在每个线程中获取到唯一实例
*/

/*
	LoadConfiguration 加载自定义配置项，
	Opt 必选配置项：
	Config_key ="name"定义配置 名，
	Config_fun ="func"定义配置 函数，
	Config_opt ="opt"	定义配置 参数选项
*/
func (a *Aurora) LoadConfiguration(options Opt) {
	a.lock.Lock()
	defer a.lock.Unlock()
	o := options()
	key, b := o[option.Config_key].(string)
	opt, b := o[option.Config_opt].(Opt)
	fun, b := o[option.Config_fun].(Configuration)
	if !b {
		//配置选项出现问题
		return
	}
	a.options[key] = &Option{
		opt,
		fun,
	}
}

// Pool 读取指定配置存入 pools中，Pool用于添加配置到池中，该方法用于初始化添加，不用于 放回池
func (a *Aurora) Pool(options Opt) {
	a.lock.Lock()
	defer a.lock.Unlock()
	opt := options()
	name := opt[option.Config_key].(string)
	o := a.options[name] //拿到对应的配置 实例
	v := o.store()
	p := &sync.Pool{}
	a.pools[name] = p
	p.Put(v)
}

// GetPool 获取容器池中的实例，前提是通过 Pool 向池中加载了。 GetPool 和 PutPool 必须成对出现
func (a *Aurora) GetPool(name string) interface{} {
	a.lock.Lock()
	defer a.lock.Unlock()
	p := a.pools[name]
	if p == nil {
		return nil
	}
	v := p.Get()
	if v == nil {
		opt := a.options[name]
		if opt == nil {
			return nil
		}
		v = opt.store()
		p.Put(v)
	}
	return v
}

// PutPool 把取出来的 实例返回放入到池中，以便下次使用
func (a *Aurora) PutPool(name string, v interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()
	p, b := a.pools[name]
	if !b {
		return
	}
	p.Put(v)
}
