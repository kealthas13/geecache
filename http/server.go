package http

import (
	"geecache/cache"
	"net/http"
)

type Server struct {
	cache.Cache//引用cache包接口
	node int//记录自己的节点
}

func (s *Server) Listen() {
	http.Handle("http://127.0.0.1", s.cacheHandler(s.node))//注册handler处理http://127.0.0.1的http协议
	http.ListenAndServe(":8088", s.cacheHandler(s.node))//监听端口
}
//New一个server并返回该对象地址
func New(c cache.Cache, i int) *Server {

	return &Server{c, i}
}
