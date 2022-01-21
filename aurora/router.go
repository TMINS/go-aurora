package aurora

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

/*
	基于字典树的路由器
	路由器规则:
		1.无法存储相同的路径
			1)形同路径的判定：校验参数相同，并且节点函数不为nil，节点函数为nil的节点说明，这个路径是未注册过，被提取为公共根
		2.路径查找按照逐层检索
		3.路由树上面存储者当前路径匹配的服务处理函数
		4.注册路径必须以 / 开头
		5.发生公共根
			1)节点和被添加路径产生公共根，提取公共根后，若公共根未注册，服务处理函数将为nil
			2)若节点恰好是公共根，则设置函数
		6.REST 风格注册
			1)同一个根路径下只能有一个REST 子路径
			2)REST 作为根路径也只能拥有一个REST 子路径
			3)REST 路径会和其它非REST同级路径发生冲突
		7.注册路径不能以/结尾（bug未修复，/user /user/ 产生 /user 的公共根 使用切割解析路径方式，解析子路径，拼接剩余子路径会存在bug ,注册路径的时候强制无法注册 / 结尾的 url）
*/

type routes interface {
	GET(string, Serve) interface{}
	POST(string, Serve) interface{}
	PUT(string, Serve) interface{}
	DELETE(string, Serve) interface{}
	HEAD(string, Serve) interface{}
}

// ServeHandler 处理函数统一实现接口
type ServeHandler interface {
	Controller(c *Ctx) interface{}
}

// Serve 注册处理函数的参数类型，它实现了 ServeHandler，处理函数调用 ServeHandler 接口时候 实际内部是调用 Serve 定义的本身
type Serve func(c *Ctx) interface{}

func (s Serve) Controller(ctx *Ctx) interface{} {
	defer func(ctx *Ctx) {
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
	}(ctx)
	return s(ctx)
}

// ServerRouter Aurora核心路由器
type route struct {
	mx          *sync.Mutex
	tree        map[string]*node // 路由树根节点
	defaultView Views            // 默认视图处理器，初始化采用 Aurora 实现的函数进行渲染
	AR          *Aurora          // Aurora 引用
}

// Node 路由节点
type node struct {
	Path      string        //当前节点路径
	handle    Handel        //服务处理函数
	Child     []*node       //子节点
	InterList []Interceptor //当前路径拦截链，默认为空
	TreeInter []Interceptor //路径匹配拦截器，默认为空
}

//——————————————————————————————————————————————————————————————————————————路由注册————————————————————————————————————————————————————————————————————————————————————————————

// addRoute 预处理被添加路径
func (r *route) addRoute(method, path string, fun Handel) {

	if path[0:1] != "/" { //校验path开头
		path += "/" + path //没有写 "/" 则添加斜杠开头
	}
	if path != "/" && path[len(path)-1:] == "/" {
		//e := fmt.Errorf(" %s 路径注册, %s 注册路径不能以 '/' 结尾", method, path)
		r.AR.auroraLog.Error(method + " registration, " + path + " the registration path cannot end with'/'")
		os.Exit(1)
	}
	if strings.HasPrefix(path, "//") { //解决 根 / 路径分组产生的bug
		path = path[1:]
	}

	if strings.Contains(path, "{}") { //此处的校验还需要加强，单一判断{}存在其他风险，开发者要么自己不能出现一些其他问题，比如 ...{}ss/.. or  .../a{s}a/.. 等情况 发现时间: 2022.1.5
		//校验注册 REST API 的路径格式。如果存在空的属性，不给注册 /asd/sad/{xx}/{xx}
		r.AR.auroraLog.Error(method + ":" + path + " The parameters of the restful interface cannot be empty {}")
		os.Exit(1)
	}
	//防止go程产生并发操作覆盖路径
	r.mx.Lock()
	defer r.mx.Unlock()
	//初始化路由树
	if r.tree == nil {
		r.tree = make(map[string]*node)
	}
	if _, ok := r.tree[method]; !ok {
		//初始化 请求类型根
		r.tree[method] = &node{}
	}
	//拿到根路径
	root := r.tree[method]

	r.add(method, root, path, fun, path) //把路径添加到根路径中中
}

// add 添加路径节点
// method 指定请求类型，root 根路径，Path和fun 被添加的路径和处理函数，path携带路径副本添加过程中不会有任何操作仅用于日志处理
// method: 请求类型(日志相关参数)
// path: 插入的路径(日志相关参数)
func (r *route) add(method string, root *node, Path string, fun Handel, path string) {

	//初始化根
	if root.Path == "" && root.Child == nil {
		root.Path = Path
		root.Child = nil
		root.handle = fun
		l := fmt.Sprintf("server interface mapping added successfully | %-10s %-20s bind to function (%-5s)", method, path, getFunName(fun))
		r.AR.auroraLog.Debug(l)
		return
	}
	if root.Path == Path { //相同路径可能是分裂或者提取的公共根
		//if root.handle != nil { //判断这个路径是否被注册过
		//	e := errors.New(method + ":" + Path + " and " + root.Path + " repeated!")
		//	monitor.En(executeInfo(e))
		//	r.AR.runtime <- monitor
		//	r.AR.initError <- e
		//} else {
		//	root.handle = fun
		//	l := fmt.Sprintf("Web Rout Mapping successds | %-10s %-20s bind (%-5s)", method, path, getFunName(fun))
		//	r.AR.message <- l
		//	return
		//}

		//此处修改，注册同样的路径，选择覆盖前一个，若是出现bug，注释掉现在使用的代码，还原上面的注释部分即可
		root.handle = fun
		l := fmt.Sprintf("server interface mapping added successfully | %-10s %-20s bind to function (%-5s)", method, path, getFunName(fun))
		r.AR.auroraLog.Debug(l)
		return
	}
	//如果当前的节点是 REST API 节点 ，子节点可以添加REST API节点
	//如果当前节点的子节点存在REST API 则不允许添加子节点

	//检擦添加路径 和 当前路径 的关系   Path:添加的路径串 path:当前root的路径（此处path只是和被添加Path区分开，并不是参数中的path）
	//1.Path 长度小于 当前path长度---> (Path 和path 有公共前缀，Path是path的父路径)
	//2.Path 长度大于 当前path长度---> (path 和Path 有公共前缀，path是path的父路径)
	//3.以上两种情况都不满足
	rootPathLength := len(root.Path)
	addPathLength := len(Path)
	if rootPathLength < addPathLength { //情况2. 当前节点可能是父节点
		if strings.HasPrefix(Path, root.Path) { //前缀检查
			i := len(root.Path)
			c := Path[i:]          //截取需要存储的 子路径
			if root.Child != nil { //若有子节点，查看是否有二级父节点
				for i := 0; i < len(root.Child); i++ {
					/*
						a:=strings.HasPrefix(root.Child[i].Path,c)
						b:=strings.HasPrefix(c,root.Child[i].Path)
						检查该节点的子节点和和要存储的子路径是否存存在父子关系
						存在父子关系则进入递归
					*/
					if strings.HasPrefix(root.Child[i].Path, c) || strings.HasPrefix(c, root.Child[i].Path) {
						r.add(method, root.Child[i], c, fun, path)
						return
					}
				}
			} else {
				//添加子节点
				if root.Child == nil {
					root.Child = make([]*node, 0)
				}
				if len(root.Child) > 0 {
					//如果存储的路径是REST API 检索 当前子节点是否存有路径，存有路径则为冲突
					for i := 0; i < len(root.Child); i++ {
						if !(strings.HasPrefix(root.Child[i].Path, "{") && strings.HasPrefix(Path, "{")) {
							//e := errors.New(method + ":" + Path + " and " + root.Child[i].Path + " collide with each other")
							r.AR.auroraLog.Error(method + ":" + Path + " and " + root.Child[i].Path + " collide with each other")
							os.Exit(1)
						}
					}
				}
				root.Child = append(root.Child, &node{Path: c, Child: nil, handle: fun})
				l := fmt.Sprintf("server interface mapping added successfully | %-10s %-20s bind to function (%-5s)", method, path, getFunName(fun))
				r.AR.auroraLog.Debug(l)
				return
			}
		}
	}
	if rootPathLength > addPathLength { //情况1.当前节点可能作为子节点
		if strings.HasPrefix(root.Path, Path) { //前缀检查
			i := len(Path)     //
			c := root.Path[i:] //需要存储的子路径，c是被分裂出来的子路径
			if root.Child != nil {
				for i := 0; i < len(root.Child); i++ {
					/*
						a:=strings.HasPrefix(root.Child[i].Path,c)
						b:=strings.HasPrefix(c,root.Child[i].Path)
						检查该节点的子节点和和要存储的子路径是否存存在父子关系
						存在父子关系则进入递归
					*/
					if strings.HasPrefix(root.Child[i].Path, c) || strings.HasPrefix(c, root.Child[i].Path) {
						r.add(method, root.Child[i], c, fun, path)
						return
					}
				}

			} else {
				//添加子节点
				if root.Child == nil {
					root.Child = make([]*node, 0)
				}
				if len(root.Child) > 0 {
					//如果存储的路径是REST API 需要检索当前子节点是否存有路径，存有路径则为冲突
					for i := 0; i < len(root.Child); i++ {
						if !(strings.HasPrefix(root.Child[i].Path, "{") && strings.HasPrefix(Path, "{")) {
							//e := errors.New(method + ":" + Path + " and " + root.Child[i].Path + " collide with each other")
							r.AR.auroraLog.Error(method + ":" + Path + " and " + root.Child[i].Path + " collide with each other")
							os.Exit(1)
						}
					}
				}
				tempChild := root.Child                                                                //保存要一起分裂的子节点
				root.Child = make([]*node, 0)                                                          //清空当前子节点  root.Child=root.Child[:0]无法清空存在bug ，直接分配保险
				root.Child = append(root.Child, &node{Path: c, Child: tempChild, handle: root.handle}) //封装被分裂的子节点 添加到当前根的子节点中
				root.Path = root.Path[:i]                                                              //修改当前节点为添加的路径
				root.handle = fun                                                                      //更改当前处理函数
				l := fmt.Sprintf("server interface mapping added successfully | %-10s %-20s bind to function (%-5s)", method, path, getFunName(fun))
				r.AR.auroraLog.Debug(l)
				return
			}
		}
	}
	//情况3.节点和被添加节点无直接关系 抽取公共前缀最为公共根
	r.merge(method, root, Path, fun, path, root.Path)
	return
}

// Merge 检测root节点 和待添加路径 是否有公共根，有则提取合并公共根
// method: 请求类型(日志相关参数)
// path: 插入的路径(日志相关参数)
// rpath: 节点的路径(日志相关参数)
// monitor: 链路日志(日志相关参数)
// root: 根合并相关参数
// Path: 根合并相关参数
// fun: 根合并相关参数
func (r *route) merge(method string, root *node, Path string, fun Handel, path string, rpath string) bool {
	pub := r.findPublicRoot(method, root.Path, Path) //公共路径
	if pub != "" {
		pl := len(pub)
		/*
			此处是提取当前节点公共根以外的操作，若当前公共根是root.Path自身则会取到空字符串 [:] 切片截取的特殊原因
			root.Path[pl:] pl是自生长度，取到最后一个字符串需要pl-1，pl取到的是个空，字符串默认为"",
			root.Path[pl:]取值为""时，说明root.Path本身就是就是公共根，只需要添加另外一个子节点到它的child切片即可
		*/
		ch1 := root.Path[pl:]
		ch2 := Path[pl:]
		if root.Child == nil {
			root.Child = make([]*node, 0)
		}
		if ch1 != "" {
			//ch1 本节点发生分裂 把处理函数也分裂 然后把当前的handler 置空,分裂的子节点也应该按照原有的顺序保留，分裂下去
			ch_child := root.Child
			root.Child = make([]*node, 0)                                                           //重新分配
			root.Child = append(root.Child, &node{Path: ch1, Child: ch_child, handle: root.handle}) //把分裂的子节点全部信息添加到公共根中
			root.handle = nil                                                                       //提取出来的公共根 没有可处理函数
		}
		if ch2 != "" {
			if len(root.Child) > 0 {
				for i := 0; i < len(root.Child); i++ {
					//单纯的被添加到此节点的子节点列表中 需要递归检测子节点和被添加节点是否有公共根
					if r.merge(method, root.Child[i], ch2, fun, path, rpath) {
						return true
					}
				}
				//检索插入路径REST API冲突。
				for i := 0; i < len(root.Child); i++ {
					if strings.HasPrefix(root.Child[i].Path, "{") || strings.HasPrefix(ch2, "{") {
						//e := errors.New(method + " :" + path + "  Conflict with other rest ful")
						r.AR.auroraLog.Error(method + " :" + path + "  Conflict with other rest ful")
						os.Exit(1)
					}
					if strings.HasPrefix(root.Child[i].Path, "{") && strings.HasPrefix(ch2, "{") {
						//e := errors.New(method + " :" + path + "  Conflict with other rest ful")
						r.AR.auroraLog.Error(method + " :" + path + "  Conflict with other rest ful")
						os.Exit(1)
					}
				}
			}
			root.Child = append(root.Child, &node{Path: ch2, Child: nil, handle: fun}) //作为新的子节点添加到当前的子节点列表中
		} else {
			//ch2为空说明 ch2是被添加路径截取的 添加的路径可能是被提出来作为公共根
			if pub == Path {
				root.handle = fun
			}
		}
		root.Path = pub //覆盖原有值设置公共根

		l := fmt.Sprintf("server interface mapping added successfully | %-10s %-20s bind to function (%-5s)", method, path, getFunName(fun))
		r.AR.auroraLog.Debug(l)
		return true
	}
	return false
}

// FindPublicRoot 查找公共前缀，如无公共前缀则返回 ""
func (r *route) findPublicRoot(method, p1, p2 string) string {
	l := len(p1)
	if l > len(p2) {
		l = len(p2) //取长度短的
	}
	index := -1
	for i := 0; i <= l && p1[:i] == p2[:i]; i++ { //此处可能发生bug
		index = i
	}
	if index > 0 && index <= l {
		//提取公共根 遇到REST API时候 需要阻止提取  主要修改 /aaa/${name} 和 /aaa/${nme} 情况下会造成提取公共根，作为rest api 参数 是不合法的
		s1 := p1[:index]
		for i := len(s1); i > 0 && s1[i-1:i] != "/"; i-- { //从后往前检索到第一个 / 如果没有遇到 $ 则没有以取 REST API为公共根
			if s1[i-1:i] == "{" {
				//fmt.Println(" REST API 路径冲突 : " + s1)
				//panic("REST API 路径冲突")
				r.AR.auroraLog.Error(method + ":" + p1 + "and" + p2 + "conflict")
				os.Exit(-1)
			}
		}
		return s1
	}
	return ""
}

// OptimizeTree 优化路由树
func (r *route) optimizeTree() {
	for _, v := range r.tree {
		optimize(v)
	}
}

// Optimize 递归排序所有子树
func optimize(root *node) {
	if root == nil {
		return
	}
	if root.Child == nil {
		return
	}
	for i := 0; i < len(root.Child); i++ {
		optimize(root.Child[i])
	}
	optimizeSort(root.Child, 0, len(root.Child)-1)
}

// OptimizeSort 对子树进行快排
func optimizeSort(child []*node, start int, end int) {
	if start < end {
		i := start
		j := end
		for i < j {
			for i < j && child[j].Path >= child[i].Path {
				j--
			}
			child[i], child[j] = child[j], child[i]
			for i < j && child[i].Path <= child[j].Path {
				i++
			}
			child[i], child[j] = child[j], child[i]
		}
		optimizeSort(child, start, i-1)
		optimizeSort(child, i+1, end)
	}
}

// ——————————————————————————————————————————————————————————————————————————————路由注册结束——————————————————————————————————————————————————————————————————————————————————

// —————————————————————————————局部拦截器——(插入拦截器查询算法和，路由查询算法一直)—————————————————————————————————————————————————————————————————————————————————————————————————

// RegisterInterceptor 向路由树上添加拦截器，添加规则只要是匹配的路径都会添加上对应的拦截器，不区分拦截的请求方式，REST API暂定还未调试支持
func (r *route) RegisterInterceptor(path string, interceptor ...Interceptor) {
	pl := len(path)
	if pl < 1 {
		//err := errors.New(path + "注册拦截器路径不能为空.")
		r.AR.auroraLog.Error("registration interceptor path cannot be empty")
		os.Exit(1)
	}
	if pl > 1 {
		//如果是通配路径 则进行校验
		if path[pl-1:] == "*" && path[pl-2:] != "/*" {
			//err := errors.New(path + " 通配符拦截器路径错误,必须以/*结尾.")
			r.AR.auroraLog.Error(path + " the wildcard interceptor path is wrong, it must end with /*")
			os.Exit(1)
		}
	}
	//校验路径开头是否以 / 否则不给添加
	if path[pl-1:] == "/" {
		//err := errors.New(path + " 拦截器不能以 / 结尾")
		r.AR.auroraLog.Error(path + " interceptor cannot end with /")
		os.Exit(1)
	}
	//为每个路径添加上拦截器
	for k, _ := range r.tree {
		r.register(r.tree[k], path, path, interceptor...)
	}
}

// register 拦截器添加和路径查找方式一样，参数path是需要添加的路径，interceptor则是需要添加的拦截器集合，root则是表示为某种请求类型进行添加拦截器
func (r *route) register(root *node, path string, lpath string, interceptor ...Interceptor) {
	if root == nil {
		return
	}
	//当前路径处理匹配成功  root.Path == path[:len(path)-2]用于匹配  /* 通配拦截器    //root.Path == path[:len(path)-1] || root.Path == path ||
	if len(path) >= 2 && root.Path == path[:len(path)-2] || root.Path == path[:len(path)-1] {
		if path[len(path)-1:] == "*" { //检测是否是通配拦截器，通配拦截器可以放在没有处理函数的路上
			//注册匹配路径
			root.TreeInter = interceptor
			if root.handle != nil {
				//添加通配拦截器 的路径是一个服务的话就一并设置
				root.InterList = interceptor
			}
			l := fmt.Sprintf("service path interceptor added successfully | %s  ", lpath)
			r.AR.auroraLog.Info(l)
			return
		}
		//root.InterList = interceptor //再次添加会覆盖通配拦截器
		//l := fmt.Sprintf("Web Rout Interceptor successds | %s ", path)
		//r.AR.message <- l
		//return
	}

	if root.Path == path && root.handle != nil {
		root.InterList = interceptor //再次添加会覆盖通配拦截器
		l := fmt.Sprintf("service path interceptor added successfully | %s  ", lpath)
		r.AR.auroraLog.Info(l)
		return
	}

	rs := strings.Split(root.Path, "/") //当前节点进行切割
	ps := strings.Split(path, "/")      //查询路径进行切割
	rsl := len(rs)
	psl := len(ps)
	sub := ""
	if psl < rsl {
		r.AR.auroraLog.Error(path + " path does not exist")
		os.Exit(1)
		//return
	}
	for i := 0; i < rsl; i++ { //解析当前路径和查找路径是否有相同部分
		//if 逐一对路径进行 比较或者解析
		if rs[i] == ps[i] { //检查rs是否和查询路径一致
			continue //如果一致则进行下一个检查
		}
		if rs[i] != ps[i] && strings.Contains(rs[i], "{") { //检测 rs是否为rest api
			if rs[i][0:1] != "{" {
				r.AR.auroraLog.Error("rest api parsing error")
				os.Exit(1)
			}
			//kl := len(rs[i])
			//key := rs[i][2 : kl-1]
			//if Args == nil {
			//	Args = make(map[string]string)
			//}
			//Args[key] = ps[i]
			continue
		} else {
			if strings.HasPrefix(ps[i], rs[i]) { //检查是否存在父子关系
				//解析被切割成为父子关系的部分
				l := len(rs[i])
				sub = ps[i][l:] //sub 被切割到子路径部分 ，子路径检索的时候需要添加到路径前面,如果sub为 "" 空则说明循环结束并没有子路径
				continue
			}
		}
		r.AR.auroraLog.Error(path + " path does not exist")
		os.Exit(1)
		//return
	}
	//此处修复 if sub=="" 为 if sub=="" && rsl==psl， /user/${name}/update  和 /user 类型情况下  /user 解析出 [""."user"]和[""."user","xxx","update"],上面的检查
	//无法检测出字串 导致/user/${name}/update 会走到/user里面，上面无法检测出子路径是查询路径和当前节点路径完全一致情况下，并且没有子路径

	if sub == "" && rsl == psl {
		//此处的路径拦截器注册暂时作用不大，后续对REST API 可能有用
		if root.handle != nil {
			if path[len(path)-1:] == "*" { //检测是否是通配拦截器
				//注册匹配路径
				root.TreeInter = interceptor
				return
			}
			root.InterList = interceptor
			return
		}
	}

	if rsl <= psl { //存在子路径  等于的情况是发生在  访问 /aa  /下面出现多个子节点 /aa是被注册需要访问的 /aa 后面没有子路径
		str := ""                    //解析子路径，用于存储下面的for循环解析的子路径
		for i := rsl; i < psl; i++ { // 检索path 剩余部分 把切割开的路径组装起来构成子路径
			if i == psl-2 && ps[psl-1] == "" { //拼接到 倒数第2个元素 判断最后一个元素为 "" 说明需要 /结尾
				str += "/" + ps[i] + "/"
				break
			}
			if i == psl-1 && ps[psl-1] != "" { //最后一个元素
				if sub != "" { //拼被丢弃的接子路径
					str += "/" + ps[i]
					str = sub + str //被丢弃的子路径是在 检索当前路径正确时候解析出来的
					break
				}
				str += "/" + ps[i]
				break
			}
			str += "/" + ps[i]
		}
		// root.Path=="/" || rsl==psl
		if rsl == psl { //遇到当前节点为 / 情况下 无法解析出 str 应为 rsl == psl，上面代码的for循环走不了， / 的子路径之下的子路径都不会以 /开头
			str = sub + str //子前缀一定要加在前面
		}
		for i := 0; i < len(root.Child); i++ { //子路径解析完成，开始遍历子节点路径，找到一个符合的路径继续走下去
			pub := r.findPublicRoot("", str, root.Child[i].Path)
			if pub != "" {
				r.register(root.Child[i], str, lpath, interceptor...)
				return
			}
		}
		//fmt.Println(path, "拦截器注册失败")
		r.AR.auroraLog.Error(path + " path does not exist")
		os.Exit(1)
		//return
	}
}

// ———————————————————————————————局部拦截器结束——————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————

// ———————————————————————————————路由查询算法—(兼职路由转发任务)——————————————————————--———————————————————————————————————————————————————————————————————————————————————————

// SearchPath 检索指定的path路由
// method 请求类型，path 查询路径，rw，req http生成的请求响应,
//ctx 和 params 主要用于转发请求时携带，初始化的时候是nil 在后续初始化，目前两个同时存在不会影响
func (r *route) SearchPath(method, path string, rw http.ResponseWriter, req *http.Request, ctx *Ctx, params HttpRequest) {
	if n, ok := r.tree[method]; ok { //查找指定的Method树
		r.search(n, path, nil, rw, req, ctx, params)
	}
}

// Search 路径查找，参数和 SearchPath意义 一致， Args map主要用于存储解析REST API路径参数，默认传nil,Interceptor拦截器可变参数，用于生成最终拦截链
func (r *route) search(root *node, path string, Args map[string]interface{}, rw http.ResponseWriter, req *http.Request, ctx *Ctx, params HttpRequest, Interceptor ...Interceptor) {

	if root == nil {
		return
	}
	if Interceptor != nil && root.TreeInter != nil {
		//把当前路径上的拦截器存起来
		for i := 0; i < len(root.TreeInter); i++ {
			Interceptor = append(Interceptor, root.TreeInter[i]) //把路径上的拦截器依次存起来
		}
	}
	//初始化参数的操作需要放在后面，狗则会导致 和 Interceptor!=nil &&  root.TreeInter!=nil 冲突重复添加一次
	if root.TreeInter != nil && Interceptor == nil {
		Interceptor = root.TreeInter
	}

	if root.Path == path { //当前路径处理匹配成功
		if root.handle != nil { //校验是否为有效路径
			//服务处理方法入口
			proxy := proxy{
				rew:             rw,
				req:             req,
				Interceptor:     true,
				HttpHandle:      root.handle,
				args:            Args,
				ctx:             ctx,
				params:          params,
				InterceptorList: root.InterList,
				TreeInter:       Interceptor,
				view:            r.defaultView, //使用路由器ServerRouter 的默认处理函数
				ar:              r.AR,
			}
			proxy.start() //开始执行
			return        //执行结束
		}
		//fmt.Println("访问路径不存在! 未注册 : " + path)
		http.NotFound(rw, req)
		return
	}

	rs := strings.Split(root.Path, "/") //当前节点进行切割
	ps := strings.Split(path, "/")      //查询路径进行切割--path是访问服务器的请求
	rsl := len(rs)
	psl := len(ps)
	sub := ""
	if psl < rsl {
		//fmt.Println("访问路径不存在! 未注册 : " + path)
		http.NotFound(rw, req)
		return
	}
	for i := 0; i < rsl; i++ { //解析当前路径和查找路径是否有相同部分
		//if 逐一对路径进行 比较或者解析
		if rs[i] == ps[i] { //检查rs是否和查询路径一致
			continue //如果一致则进行下一个检查
		}
		if rs[i] != ps[i] && strings.Contains(rs[i], "{") { //检测 rs是否为rest api
			if rs[i][0:1] != "{" && rs[i][len(rs[i])-1:] != "}" { // rs[i][0:1] != "{"  添加了修改了 参数解析检查
				panic("REST API 解析错误")
			}
			kl := len(rs[i])
			key := rs[i][1 : kl-1]
			if Args == nil {
				Args = make(map[string]interface{})
			}
			Args[key] = ps[i]
			continue
		} else {
			if strings.HasPrefix(ps[i], rs[i]) { //检查是否存在父子关系
				//解析被切割成为父子关系的部分
				l := len(rs[i])
				sub = ps[i][l:] //sub 被切割到子路径部分 ，子路径检索的时候需要添加到路径前面,如果sub为 "" 空则说明循环结束并没有子路径
				continue
			}
		}
		//fmt.Println("访问路径不存在! 未注册 : " + path)
		http.NotFound(rw, req)
		return
	}
	//此处修复 if sub=="" 为 if sub=="" && rsl==psl， /user/${name}/update  和 /user 类型情况下  /user 解析出 [""."user"],[""."user","xxx","update"],上面的检查
	//无法检测出字串 导致/user/${name}/update 会走到/user里面，上面无法检测出子路径是查询路径和当前节点路径完全一致情况下，并且没有子路径
	if sub == "" && rsl == psl {
		if root.handle != nil {
			//服务处理方法入口
			proxy := proxy{
				rew:             rw,
				req:             req,
				Interceptor:     true,
				HttpHandle:      root.handle,
				args:            Args,
				ctx:             ctx,
				InterceptorList: root.InterList,
				view:            r.defaultView,
				ar:              r.AR,
				TreeInter:       Interceptor,
			}
			proxy.start()
			return
		}
	}
	if rsl <= psl { //存在子路径  等于的情况是发生在  访问 /aa  /下面出现多个子节点 /aa是被注册需要访问的 /aa 后面没有子路径
		str := ""                    //解析子路径，用于存储下面的for循环解析的子路径
		for i := rsl; i < psl; i++ { // 检索path 剩余部分 把切割开的路径组装起来构成子路径
			if i == psl-2 && ps[psl-1] == "" { //拼接到 倒数第2个元素 判断最后一个元素为 "" 说明需要 /结尾
				str += "/" + ps[i] + "/"
				break
			}
			if i == psl-1 && ps[psl-1] != "" { //最后一个元素
				if sub != "" { //拼被丢弃的接子路径
					str += "/" + ps[i]
					str = sub + str //被丢弃的子路径是在 检索当前路径正确时候解析出来的
					break
				}
				str += "/" + ps[i]
				break
			}
			str += "/" + ps[i]
		}
		// root.Path=="/" || rsl==psl
		if rsl == psl { //遇到当前节点为 / 情况下 无法解析出 str 应为 rsl == psl，上面代码的for循环走不了， / 的子路径之下的子路径都不会以 /开头
			str = sub + str //子前缀一定要加在前面
		}
		for i := 0; i < len(root.Child); i++ { //子路径解析完成，开始遍历子节点路径，找到一个符合的路径继续走下去
			pub := r.findPublicRoot("", str, root.Child[i].Path)
			if pub != "" {
				r.search(root.Child[i], str, Args, rw, req, ctx, params, Interceptor...)
				return
			}
		}
		//fmt.Println("访问路径不存在! 未注册 : " + path)
		http.NotFound(rw, req)
		return
	}
}

// ———————————————————————————————路由查询算法结束——————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————

// ServeHTTP 一切的开始
func (a *Aurora) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	mapping := req.URL.Path
	if err := checkUrl(mapping); err != nil {
		rw.Header().Set(contentType, a.resourceMapType[".json"])
		http.Error(rw, newErrorResponse(mapping, "rest ful 路径格式不正确，不能包含符号'{'或'}',"+err.Error(), 500), 500)
		return
	}

	if index := strings.LastIndex(mapping, "."); index != -1 { //此处判断这个请求可能为静态资源处理
		t := mapping[index:] //截取可能的资源类型
		a.resourceHandler(rw, req, mapping, t)
		return
	}

	a.router.SearchPath(req.Method, req.URL.Path, rw, req, nil, nil) //初始一个nil ctx
}

func getFunName(fun Handel) string {
	return runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
}

func checkUrl(url string) error {
	if strings.Contains(url, "%7B") || strings.Contains(url, "%7D") {
		return errors.New("url请求是非法的")
	}
	return nil
}

func restCheck(url string) bool {
	if strings.Contains(url, "{") && strings.Contains(url, "}") {
		ulen := len(url)
		ubyte := []byte(url)
		for i := 0; i < ulen; i++ {
			if string(ubyte[i]) == "{" {
				//带实现校验
			}
		}
	}
	return true
}
