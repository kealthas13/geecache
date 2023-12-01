package main

import (
	"flag"
	"fmt"
	"geecache/cache"
	"geecache/http"
	"strings"
)

type interval []string

// 实现String()方法
func (h *interval) String() string {
	return fmt.Sprintf("%v", *h)
}

// Set 方法
func (h *interval) Set(s string) error {
	for _, v := range strings.Split(s, ",") {
		*h = append(*h, v)
	}
	return nil
}

var (
	node int
)

func init() {
	flag.IntVar(&node, "n", 0, "node")
}

func main() {
	flag.Parse() //读入参数-n（node）

	if node > 3 || node < 1 {
		fmt.Println("the serial number of node must be 1-3")
		return
	}
	//创建一块在内存中的cache
	c := cache.New()
	//以cache、node为参数创建一个指向http.server的结构体
	s := http.New(c, node)
	//以cache、node为参数创建一个gRPC服务端并运行
	http.New_RPC(node, c)
	//监听localhost:8088
	s.Listen()

}
