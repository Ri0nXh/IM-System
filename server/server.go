package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// 定义一个server结构体
/*
OnlineMap 存储一个全局在线的用户map
Msg 是一个接收消息并实现全局广播
*/
type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	Lock      sync.RWMutex
	Msg       chan string
}

// 初始化server
func NewServer(ip string, port int) *Server {
	s := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Msg:       make(chan string),
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
	go s.ListenMessager()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept conn error: ", err)
			continue
		}
		// 处理用户业务
		go s.Handler(conn)
	}
}

// 用户业务处理
func (s *Server) Handler(conn net.Conn) {
	// 初始化用户信息，并加入onlinemap中
	user := NewUser(conn, s)

	// 调用广播方法，去发送消息。
	//s.BroadCast(user, "is online, come chat!!!")
	user.Online()

	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "is offline...")
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("receive client msg error:", err)
				continue
			}

			// 去除最后的\n字符
			revicedMsg := string(buf[:n-1])

			// 将用户发送的消息进行广播。（感觉有点像群聊）
			//s.BroadCast(user, revicedMsg)
			user.DoMessage(revicedMsg)
		}
	}()
	select {}
}

// 广播消息发送者
func (s *Server) BroadCast(u *User, msg string) {
	sendMsg := fmt.Sprintf("[%s] %s\n", u.Addr, msg)
	s.Msg <- sendMsg
}

// 时刻监听服务端发送过来的消息（消息接收者）
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Msg
		s.Lock.Lock()
		for _, userInfo := range s.OnlineMap {
			userInfo.Msg <- msg
		}
		s.Lock.Unlock()
	}
}
