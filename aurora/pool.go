package aurora

import (
	"github.com/awensir/go-aurora/aurora/option"
	"sync"
)

/*
	通过pool的方式添加 配置，能够保证在每个线程中获取到唯一实例
*/

// Pool 读取指定配置存入 pools中，Pool用于添加配置到池中，该方法用于初始化添加，不用于 放回池
func (a *Aurora) Pool(options Opt) {
	a.lock.Lock()
	defer a.lock.Unlock()
	opt := options()
	name := opt[option.Config_key].(string)
	o, b := a.config[name] //拿到对应的配置 实例
	if o == nil && !b {
		// 不存在对应配置
		return
	}
	v := o.store() //加载并 得到配置实例本身
	p := &sync.Pool{}
	a.pools[name] = p //初始化改 name的池
	p.Put(v)          //变量放入池中
}

// GetPool 获取容器池中的实例，前提是通过 Pool 向池中加载了。 GetPool 和 PutPool 必须成对出现
func (a *Aurora) GetPool(name string) interface{} {
	a.lock.Lock()
	defer a.lock.Unlock()
	p := a.pools[name]
	if p == nil {
		// name 池不存在 或者未初始化
		return nil
	}
	v := p.Get()
	if v == nil {
		// 如果池中的变量自行销毁 则使用配置重新 初始化
		opt := a.config[name]
		if opt == nil {
			return nil
		}
		v = opt.store()
		//p.Put(v)  初始化后不放入池中,后续在调用 PutPool(name,v) 中放回
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
