package router

import (
	"fmt"
	"strings"
	"testing"
)

func T1() {
	fmt.Println("t1")
}
func T2() {
	fmt.Println("t2")
}
func T3() {
	fmt.Println("t3")
}
func T4() {
	fmt.Println("t4")
}
func T5() {
	fmt.Println("t5")
}
func TestLen(t *testing.T) {
	s:="/abcd"
	s2:="/abcd/b/we"
	arr:=strings.Split(s,"/")
	arr2:=strings.Split(s2,"/")
	t.Log(arr,arr2)
	
	str:=""
	for i:= len(arr);i< len(arr2);i++ {
		if i==i-2 && arr2[i-1]=="" { //拼接到 倒数第2个元素 判断最后一个元素为 "" 说明需要 /结尾
			str+="/"+arr2[i]+"/"
			break
		}
		if i==i-1 && arr2[i-1]!="" { //最后一个元素
			str+=arr2[i]
			break
		}
		str+="/"+arr2[i]
	}
	
	arrl:=len(arr)
	for i:=0;i< arrl;i++{
		if i==arrl-2 && arr[arrl-1]=="" { //最后一个元素为 "" 说明需要 /结尾
			str+=arr[i]+"/"
			break
		}
		if i==arrl-1 && arr[arrl-1]!="" {
			str+=arr[i]
			break
		}
		str+=arr[i]+"/"
	}
	
	
}

func TestRouter(t *testing.T) {
	r:=Router{}
	r.addRoute("GET","/abcd/b/we",T1)
	r.addRoute("GET","/ab/a/we",T2)
	r.addRoute("GET","/ab/a",T5)
	r.addRoute("GET","/abcd/c/we",T3)
	r.addRoute("GET","/abcd/a/we",T4)
	//r.OptimizeTree()
	r.SearchPath("GET","/abcd/b/we")
	r.SearchPath("GET","/ab/a/we")
	r.SearchPath("GET","/abcd/c/we")
	r.SearchPath("GET","/abcd/a/we")
	t.Log(r)
}

func TestFindRouter(t *testing.T) {
	r:=Router{}
	r.addRoute("GET","/abc/b/we",T1)
	r.addRoute("GET","/ab/a/we",T2)
	r.addRoute("GET","/ab/a",T5)
	r.addRoute("GET","/abcdd/c/we",T3)
	r.addRoute("GET","/abcdmb/a/we",T4)
	r.addRoute("GET","/abcd/${name}/ww/${age}/ca",T1)
	r.addRoute("GET","/abcd/${name}/ww/${age}/cab/bcc",T2)
	r.addRoute("GET","/abcd/${name}/ww/${age}/ca/bcace",T3)
	r.OptimizeTree()
	r.SearchPath("GET","/abc/b/we")
	r.SearchPath("GET","/ab/a")
	r.SearchPath("GET","/abcdd/c/we")
	r.SearchPath("GET","/abcdmb/a/we")
	r.SearchPath("GET","/abcd/a/ww/11/ca")
	t.Log(r)
}

func TestString(t *testing.T) {
	a:="/"
	arr:=strings.Split(a,"/")
	t.Log(arr)
}

func TestString2(t *testing.T) {
	p:="/ab/a/we"
	t.Log(p[:])
}

func TestSort(t *testing.T) {
	arr:=[]int{0,1,5,2,0,7,8,6,7}
	sort(arr,0,len(arr)-1)
	t.Log(arr)
	f:=finds(arr,0)
	t.Log(f)
}

func finds(a []int,n int) bool {
	if len(a)<=0{
		return false
	}
	//计算中轴
	l:= len(a)/2
	if l==0{
		if n!=a[l] {
			return false
		}
	}
	if n==a[l]{
		return true
	}
	if l>0 && n<a[l]{
		//查找的数字小于中轴 递归数组左半部分
		return finds(a[:l],n)
	}
	if  n>a[l] {
		return finds(a[l:],n)
	}
	return false
}
var c int=0
func sort(a []int,start int,end int)  {
	if start < end {
		
		i:=start
		j:=end
		//k:=a[start]
		for i<j {
			for i<j && a[j]>=a[i] {
				j--
			}
			//k,a[j]=a[j],k
			a[j],a[i]=a[i],a[j]
			for i<j && a[i]<=a[j] {
				i++
			}
			//k,a[i]=a[i],k
			a[j],a[i]=a[i],a[j]
		}
		c++
		fmt.Println("第: ",c,"次")
		fmt.Println(a)
		sort(a,start,i-1)
		sort(a,i+1,end)
	}
}
