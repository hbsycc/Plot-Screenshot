package demo

import (
	"fmt"
	"time"
)

// 无参函数
func f1() {
	fmt.Println("--> in func f1")
	fmt.Println("-->(f1) f1 node")
}

// 调用例子
func goFunc1(f func()) {
	go f()
}

// 有定长参数的函数
func f2(i interface{}) {
	fmt.Println("--> in func f2")
	fmt.Println("-->(f2) f2 node", i)
}

// 调用例子
func goFunc2(f func(interface{}), i interface{}) {
	go f(i)
}

// 有变长参数的函数
func f3(args ...interface{}) {
	fmt.Println("--> in func f3")
	fmt.Println("-->(f3) f3 done", args)
}

// 复杂调用例子
func goFunc3(f interface{}, args ...interface{}) {
	if len(args) > 1 {
		go f.(func(...interface{}))(args)
	} else if len(args) == 1 {
		go f.(func(interface{}))(args[0])
	} else {
		go f.(func())()
	}
}

func main() {
	//-------- 演示，函数作为参数传递
	goFunc1(f1)      //无参数
	goFunc2(f2, 100) //有参数

	//复杂用法
	goFunc3(f1)
	goFunc3(f2, "xxx")
	goFunc3(f3, "hello", "world", 1, 3.14)

	// ------------- 程序延时退出------------
	time.Sleep(5 * time.Second)
}
