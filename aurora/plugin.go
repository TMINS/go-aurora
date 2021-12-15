package aurora

type PluginFunc func(ctx *Ctx)

// Plugin 加载全局插件
func (a *Aurora) Plugin(plugs ...PluginFunc) {
	if a.plugins == nil {
		a.plugins = make([]PluginFunc, 0)
		for _, v := range plugs {
			a.plugins = append(a.plugins, v)
		}
		return
	}
	for _, v := range plugs {
		a.plugins = append(a.plugins, v)
	}
}
