package http

import (
	"context"
	"fmt"
	"geecache/cache"
	pb "geecache/proto"
	"net"
	"net/http"

	"google.golang.org/grpc"
)
//其他节点的IP
const network01 string = "188.168.0.101:"
const network02 string = "188.168.0.102:"
const network03 string = "188.168.0.103:"

type server_rpc struct {
	pb.UnimplementedNodeServiceServer
	cache.Cache
}

// 服务端
func (s *server_rpc) SendAndResponse(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	var responses int
	var key, value string
	if req.Action == http.MethodPost {//POST实现

		v := []byte(req.HttpValue) //转成byte
		if len(req.HttpValue) != 0 {//如果value不空
			e := s.Set(req.KeyValue, v) //存入
			if e != nil {//异常
				// log.Println(e)
				responses = http.StatusInternalServerError
			} else {//成功
				responses = http.StatusOK
			}
		}
	} else if req.Action == http.MethodGet {
		//GET实现
		b, e := s.Get(req.KeyValue)
		// fmt.Println(b, e)

		if e != nil {//异常
			// log.Println(e)
			responses = http.StatusInternalServerError
			// return
		} else if len(b) == 0 {//没找到
			responses = http.StatusNotFound
			// return
		} else {//成功
			value = string(b)
			responses = http.StatusOK
		}

	} else if req.Action == http.MethodDelete {
		//DEL实现
		b, _ := s.Get(req.KeyValue)
		responses = http.StatusOK//永远返回HTTP 200
		if len(b) != 0 {//存在
			value = "1"
			_ = s.Del(req.KeyValue)
		} else {//不存在
			value = "0"
		}
	}
	return &pb.Response{Action: req.Action, KeyValue: key, HttpValue: value, Response: int64(responses)}, nil//调用返回
}
//创建gRPC服务端
func New_RPC(node int, c cache.Cache) {
	var network string
	//选择本地ip
	switch node {
	case 1:
		network = network01
		break
	case 2:
		network = network02
		break
	case 3:
		network = network03
		break
	default:
		fmt.Println("the serial number of node must be 1-3")
		break
	}
	//开启端口
	listen, err := net.Listen("tcp", network+"12345")
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}
	//创建rpc服务
	grpcServer := *grpc.NewServer()
	s := &server_rpc{}
	//把cache接口引入
	s.Cache = c
	//在服务端注册编写的服务
	pb.RegisterNodeServiceServer(&grpcServer, s)

	fmt.Println("server listening at", listen.Addr())
	//启用一个子进程启动服务
	go func() {
		err = grpcServer.Serve(listen)
		if err != nil {
			fmt.Println("error")
		}
	}()

}
