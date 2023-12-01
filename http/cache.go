package http

import (
	"context"
	"encoding/json"
	"fmt"
	pb "geecache/proto"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 实现http.Handler接口
type cacheHandler struct {
	*Server
	node int
}

// const com int = 9527

// gRPC客户端,用于连接寻找其他节点服务端，用完即删
func connect(key string, value string, h cacheHandler, Method string, i int) (string, string, string, int) {
	var com, network string //选择网络

	switch i {
	case 1:
		network = network01
		com = com1
		break
	case 2:
		network = network02
		com = com2
		break
	case 3:
		network = network03
		com = com3
		break
	default:
		fmt.Println("the serial number of node must be 1-3")
		break
	}
	//连接到server
	connect, err := grpc.Dial(network+com, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	//返回前抛掉
	defer connect.Close()
	//建立连接
	client := pb.NewNodeServiceClient(connect)
	//执行rpc调用
	resp, _ := client.SendAndResponse(context.Background(), &pb.Request{Action: Method, KeyValue: key, HttpValue: value})
	//返回服务端取回数据
	return resp.Action, resp.KeyValue, resp.HttpValue, int(resp.Response)

}

// 实现http.handle接口
func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//约定http API
	r.Header.Set("Content-Type", "application/json;charset=UTF-8")
	// w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	keys := strings.Split(r.URL.Path, "/") //按/分割字符

	key := keys[len(keys)-1] //取出最后一个，主要用于get方法

	// if len(key) == 0 { //简单判断防止空取
	// 	fmt.Println("没有有效key")
	// 	return
	// }

	b, _ := ioutil.ReadAll(r.Body) //读出body
	s := string(b)                 //转译string格式
	var value string               //定义k-v的v
	if len(s) != 0 {               //减少不必要操作
		value = s[strings.IndexByte(s, ':')+2 : strings.IndexByte(s, '}')] //切片取value
	}

	words := make(map[string]*string) //定义字符串map
	json.Unmarshal(b, &words)         //body反序列化成字符串map

	for k := range words {
		key = k //取出key关键词，此时因为只有一个kv，多了这样不行，考虑到作业性质不做容错了
	}
	hash_i := String_hash(key) //hash操作目标节点
	if hash_i == h.node {      //如果在本地
		h.Method_deal(key, value, w, r) //调用本地的处理方法
		return
	} else { //如果不是本地，调用封装gRPC的connect（）
		_, _, HttpValue, Response := connect(key, value, *h, r.Method, hash_i)
		if r.Method == http.MethodPost { //如果是post，直接返回（可以改进判断成功与否）
			return
		}
		w.WriteHeader(int(Response))    //将返回的http封入
		if r.Method == http.MethodGet { //如果是get

			if Response == http.StatusOK { //get成功
				s := "{" + "\"" + key + "\"" + ":" + " " + HttpValue + "}" //处理输出格式
				str := []byte(s)                                           //转译成byte
				w.Write(str)                                               //封入body体
				return
			}
			return

		}
		if r.Method == http.MethodDelete { //删除方法
			str := []byte(HttpValue) //删除数量
			w.Write(str)             //封入body体
			return
		}

	}
	return
}

// 本地的set、get、del处理方法
func (h *cacheHandler) Method_deal(key string, value string, w http.ResponseWriter, r *http.Request) {

	m := r.Method //请求方法

	if m == http.MethodPost { //POST实现

		v := []byte(value) //转成byte
		if len(value) != 0 {
			e := h.Set(key, v) //存入
			if e != nil {
				// log.Println(e)
				w.WriteHeader(http.StatusInternalServerError) //存入失败
			} else {
				w.WriteHeader(http.StatusOK) //成功
			}
		}
		return
	} else if m == http.MethodGet {
		//GET实现
		b, _ := h.Get(key)
		//删除了一些考虑的异常情况
		// fmt.Println(b, e)
		// if e != nil {
		// 	log.Println(e)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		if len(b) == 0 { //找不到
			w.WriteHeader(http.StatusNotFound) //404
			// w.Write(b1)
			return
		}
		s := "{" + "\"" + key + "\"" + ":" + " " + string(b) + "}" //顺利则处理输出格式
		b1 := []byte(s)
		w.Write(b1) //写入body体
		return
	} else if m == http.MethodDelete {
		//DEL实现
		b, e := h.Get(key)
		//删除了一些考虑的异常情况
		// if e != nil {
		// 	log.Println(e)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		w.WriteHeader(http.StatusOK) //永远返回HTTP 200
		if len(b) == 0 {             //如果not found
			str := "0"
			body := []byte(str)
			w.Write(body)
			return
		}
		e = h.Del(key)
		if e == nil { //删除成功
			// log.Println(e)
			str := "1"
			body := []byte(str)
			w.Write(body)
		}
		return
	}
	// w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) cacheHandler(i int) http.Handler {
	return &cacheHandler{s, i}
}
