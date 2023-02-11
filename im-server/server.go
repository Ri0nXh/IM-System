package im_server

import (
	"fmt"
	"net"
)

// 定义一个server结构体
type Server struct {
	Ip   string
	Port int
}

// 初始化server
func NewServer(ip string, port int) *Server {
	s := &Server{
		Ip:   ip,
		Port: port,
	}
	return s
}

// 启动server
func (s *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("conn tcp im-server error：", err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept conn error: ", err)
			continue
		}
		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	fmt.Printf("[%s] 连接进来了，处理业务", conn.RemoteAddr().String())
}
