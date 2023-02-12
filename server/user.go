package server

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	Conn   net.Conn
	Msg    chan string
	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	u := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		Conn:   conn,
		Msg:    make(chan string),
		Server: server,
	}
	// 监听消息
	go u.ListenerMsg()
	return u
}

// 用户上线方法
func (u *User) Online() {
	u.Server.Lock.Lock()
	u.Server.OnlineMap[u.Name] = u
	u.Server.Lock.Unlock()

	u.Server.BroadCast(u, "online ! ")
}

// 用户下线方法
func (u *User) Offline() {
	u.Server.Lock.Lock()
	delete(u.Server.OnlineMap, u.Name)
	u.Server.Lock.Unlock()

	u.Server.BroadCast(u, "offline ! ")
}

func (u *User) DoMessage(msg string) {
	u.Server.BroadCast(u, msg)
}
func (u *User) ListenerMsg() {
	for {
		msg := <-u.Msg
		u.Conn.Write([]byte(msg))
	}
}
