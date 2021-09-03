package router

import (
	"strings"
)
/*
	路由存储规则参考HttpRouter
	基于查询树的路由器
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
			1)同一个根路径下只能有一个REST 子路径 （***）
			2)REST 作为根路径也只能拥有一个REST 子路径
			3)REST 路径会和其它非REST同级路径发生冲突
*/

const (
	REST_API = "(\\$){1}(\\{){1}([A-Za-z0-9])*(\\}){1}"
	REST_ARGS = "[A-Za-z0-9]*"
)

type HandlerFun func()
type Node struct {
	Path string         //当前路径
	handle HandlerFun
	Child []*Node       //子节点
}

type Router struct {
	tree map[string]*Node
}

// addRoute 预处理被添加路径
func (r *Router) addRoute(method,path string,fun HandlerFun)  {
	if r.tree == nil {
		r.tree=make(map[string]*Node)       //初始化路由树
	}
	if _,ok:=r.tree[method];!ok{
		r.tree[method]=&Node{}              //初始化制定
	}
	root:=r.tree[method]                    //拿到根路径
	if path[0:1]!="/"{                      //校验path开头
		path+="/"+path                      //没有写 "/" 则添加斜杠开头
	}
	
	r.add(root,path,fun)                    //把路径添加到根路径中中
}

func (r *Router) add(root *Node,Path string,fun HandlerFun)  {
	
	//初始化根
	if root.Path==""&&root.Child==nil {
		root.Path=Path
		root.Child=nil
		root.handle=fun
		return
	}
	if root.Path==Path{                //相同路径可能是分裂或者提取的公共根
		if root.handle!=nil {         //判断这个路径是否被注册过
			panic(Path+" repetition ")
		}else {
			root.handle=fun           //
		}
	}
	//如果当前的节点是 REST API 节点 ，子节点可以添加REST API节点
	//如果当前节点的子节点存在REST API 则不允许添加子节点
	
	
	//检擦添加路径 和 当前路径 的关系   Path:添加的路径串 path:当前root的路径
	//1.Path 长度小于 当前path长度---> (Path 和path 有公共前缀，Path是path的父路径)
	//2.Path 长度大于 当前path长度---> (path 和Path 有公共前缀，path是path的父路径)
	//3.以上两种情况都不满足
	rootPathLength:= len(root.Path)
	addPathLength:= len(Path)
	if rootPathLength < addPathLength {         //情况2. 当前节点可能是父节点
		if strings.HasPrefix(Path,root.Path) {  //前缀检查
			i:= len(root.Path)
			c:=Path[i:]             //截取需要存储的 子路径
			if root.Child!=nil {    //若有子节点，查看是否有二级父节点
				for i:=0;i< len(root.Child);i++ {
					/*
						a:=strings.HasPrefix(root.Child[i].Path,c)
						b:=strings.HasPrefix(c,root.Child[i].Path)
						检查该节点的子节点和和要存储的子路径是否存存在父子关系
						存在父子关系则进入递归
					*/
					if strings.HasPrefix(root.Child[i].Path,c) || strings.HasPrefix(c,root.Child[i].Path) {
						r.add(root.Child[i],c,fun)
						return
					}
				}
			}else {
				//添加子节点
				if root.Child==nil{
					root.Child=make([]*Node,0)
				}
				if len(root.Child)>0{
					//如果存储的路径是REST API 检索 当前子节点是否存有路径，存有路径则为冲突
					for i:=0;i< len(root.Child);i++{
						if !(strings.HasPrefix(root.Child[i].Path,"$") && strings.HasPrefix(Path,"$")){
							panic(Path+" conflicting "+root.Child[i].Path)
						}
					}
				}
				root.Child=append(root.Child,&Node{Path: c,Child: nil,handle: fun})
				return
			}
		}
	}
	if rootPathLength > addPathLength {  //情况1.当前节点可能作为子节点
		if strings.HasPrefix(root.Path,Path) {  //前缀检查
			i:= len(Path)            //
			c:=root.Path[i:]         //需要存储的子路径，c是被分裂出来的子路径
			if root.Child!=nil {
				for i:=0;i< len(root.Child);i++ {
					/*
						a:=strings.HasPrefix(root.Child[i].Path,c)
						b:=strings.HasPrefix(c,root.Child[i].Path)
						检查该节点的子节点和和要存储的子路径是否存存在父子关系
						存在父子关系则进入递归
					*/
					if strings.HasPrefix(root.Child[i].Path,c) || strings.HasPrefix(c,root.Child[i].Path) {
						r.add(root.Child[i],c,fun)
						return
					}
				}
				
			}else {
				//添加子节点
				if root.Child==nil{
					root.Child=make([]*Node,0)
				}
				if len(root.Child)>0{
					//如果存储的路径是REST API 需要检索当前子节点是否存有路径，存有路径则为冲突
					for i:=0;i< len(root.Child);i++{
						if !(strings.HasPrefix(root.Child[i].Path,"$") && strings.HasPrefix(Path,"$")){
							panic(Path+" conflicting "+root.Child[i].Path)
						}
					}
				}
				tempchild:=root.Child   //保存要一起分裂的子节点
				root.Child=make([]*Node,0) //清空当前子节点  root.Child=root.Child[:0]无法清空存在bug ，直接分配保险
				root.Child=append(root.Child,&Node{Path: c,Child: tempchild,handle: root.handle}) //封装被分裂的子节点 添加到当前根的子节点中
				root.Path=root.Path[:i]                   //修改当前节点为添加的路径
				root.handle=fun                           //更改当前处理函数
				return
			}
		}
	}
	//情况3.节点和被添加节点无直接关系 抽取公共前缀最为公共根
	Merge(root,Path,fun)
	return
}

// Merge 检测root节点 和待添加路径 是否有公共根，有则提取合并公共根
func Merge(root *Node,Path string,fun HandlerFun) bool {
	pub:=FindPublicRoot(root.Path,Path)  //公共路径
	if pub !="" {
		pl:= len(pub)
		/*
			此处是提取当前节点公共根以外的操作，若当前公共根是root.Path自身则会取到空字符串 [:] 切片截取的特殊原因
			root.Path[pl:] pl是自生长度，取到最后一个字符串需要pl-1，pl取到的是个空，字符串默认为"",
			root.Path[pl:]取值为""时，说明root.Path本身就是就是公共根，只需要添加另外一个子节点到它的child切片即可
		*/
		ch1:=root.Path[pl:]
		ch2:=Path[pl:]
		if root.Child==nil {
			root.Child=make([]*Node,0)
		}
		if ch1!="" {
			//ch1 本节点发生分裂 把处理函数也分裂 然后把当前的handler 置空,分裂的子节点也应该按照原有的顺序保留，分裂下去
			ch_child:=root.Child
			root.Child=make([]*Node,0) //重新分配
			root.Child=append(root.Child,&Node{Path: ch1,Child: ch_child,handle: root.handle}) //把分裂的子节点全部信息添加到公共根中
			root.handle=nil  //提取出来的公共根 没有可处理函数
		}
		if ch2!=""{
			if len(root.Child)>0 {
				for i:=0;i< len(root.Child);i++{
					//单纯的被添加到此节点的子节点列表中 需要递归检测子节点和被添加节点是否有公共根
					if Merge(root.Child[i],ch2,fun) {
						return true
					}
				}
				//检索插入路径REST API冲突
				for i:=0;i< len(root.Child);i++{
					if strings.HasPrefix(root.Child[i].Path,"$") || strings.HasPrefix(ch2,"$"){
						panic(Path+"  conflicting  "+root.Child[i].Path)
					}
					if strings.HasPrefix(root.Child[i].Path,"$") && strings.HasPrefix(ch2,"$"){
						panic(Path+"  conflicting  "+root.Child[i].Path)
					}
				}
			}
			root.Child=append(root.Child,&Node{Path: ch2,Child: nil,handle: fun}) //作为新的子节点添加到当前的子节点列表中
		}else {
			//ch2为空说明ch2 是被添加路径截取的 被添加路径可能是被提出来作为公共根
			if pub==Path{
				root.handle=fun
			}
		}
		root.Path=pub  //覆盖原有值设置公共根
		return true
	}
	return false
}

// FindPublicRoot 查找公共前缀，如无公共前缀则返回 ""
func FindPublicRoot(p1, p2 string) string {
	l:= len(p1)
	if l> len(p2){
		l= len(p2)  //取长度短的作比较
	}
	index:=-1
	for i:=0;i<=l&&p1[:i]==p2[:i];i++{    //此处可能发生bug
		index=i
	}
	if index>0 && index<=l {
		//提取公共根 遇到REST API时候 需要阻止提取  主要修改 /abcd/${name} 和 /abcd/${nme} 情况下会造成提取公共根，作为rest api 参数 是不合法的
		s1:=p1[:index]
		for i:= len(s1);i>0 && s1[i-1:i]!="/";i-- { //从后往前检索到第一个 / 如果没有遇到 $ 则没有以取 REST API为公共根
			if s1[i-1:i]=="$"{
				panic("REST API 路径冲突")
			}
		}
		return s1
	}
	return ""
}

// OptimizeTree 优化路由树
func (r Router) OptimizeTree() {
	 for _,v:=range r.tree{
	 	Optimize(v)
	 }
}

// Optimize 递归排序所有子树
func Optimize(root *Node)  {
	if root==nil{
		return
	}
	if root.Child==nil{
		return
	}
	for i:=0;i< len(root.Child);i++{
		Optimize(root.Child[i])
	}
	OptimizeSort(root.Child,0, len(root.Child)-1)
}

// OptimizeSort 对子树进行快排
func OptimizeSort(child []*Node, start int, end int)  {
	if start<end{
		i:=start
		j:=end
		for i<j {
			for i<j && child[j].Path>=child[i].Path {
				j--
			}
			child[i],child[j]=child[j],child[i]
			for i<j && child[i].Path<=child[j].Path {
				i++
			}
			child[i],child[j]=child[j],child[i]
		}
		OptimizeSort(child,start,i-1)
		OptimizeSort(child,i+1,end)
	}
}


func (r Router) SearchPath(method,path string) {
	if n,ok:=r.tree[method];ok{
		r.Search(n,path,nil)
	}
}

func (r Router) Search(root *Node,path string,Args map[string]string) {
	if root==nil{
		return
	}
	rs:=strings.Split(root.Path,"/") //当前节点进行切割
	ps:=strings.Split(path,"/")      //查询路径进行切割
	rsl:= len(rs)
	psl:= len(ps)
	sub:=""
	for i:=0;i<rsl;i++{     //解析当前路径和查找路径是否有相同部分
		//if 逐一对路径进行 比较或者解析
		if rs[i]==ps[i] {   //检查rs是否和查询路径一致
			continue        //如果一致则进行下一个检查
		}
		if rs[i]!=ps[i] && strings.Contains(rs[i],"$"){ //检测 rs是否为rest api
			if rs[i][0:1]!="$"{
				panic("REST API 解析错误")
			}
			kl:= len(rs[i])
			key:=rs[i][2:kl-1]
			if Args==nil{
				Args=make(map[string]string)
			}
			Args[key]=ps[i]
			continue
		}else {
			if strings.HasPrefix(ps[i],rs[i]){ //检查是否存在父子关系
				//解析被切割成为父子关系的部分
				l:= len(rs[i])
				sub=ps[i][l:]
				continue
			}
		}
		panic("访问路径不存在! 未注册")
	}
	//当前路径处理没问题，查看当前路径和查询路径的长度，
	if rsl == psl {
		if root.handle!=nil {
			root.handle()       //服务处理方法入口
			return
		}
		panic("访问路径不存在! 未注册")
	}
	if rsl < psl {   //存在子路径
		str:=""
		//说明存在子节点要处理
		for i:=rsl;i<psl;i++ {
			if i==psl-2 && ps[psl-1]=="" { //拼接到 倒数第2个元素 判断最后一个元素为 "" 说明需要 /结尾
				str+="/"+ps[i]+"/"
				break
			}
			if i==psl-1 && ps[psl-1]!="" { //最后一个元素
				if sub!=""{                //拼接子路径
					str+="/"+ps[i]
					str=sub+str
					break
				}
				str+="/"+ps[i]
				break
			}
			str+="/"+ps[i]
		}
		
		for i:=0;i< len(root.Child);i++ {
			pub:=FindPublicRoot(str,root.Child[i].Path)
			if pub!="" {
				r.Search(root.Child[i],str,Args)
				return
			}
		}
	}
}



