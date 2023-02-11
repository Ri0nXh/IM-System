package im_user

import (
	"net"
)

type User struct {
	Name string
	Addr string
	Conn net.Conn
	Msg  chan string
}

func NewUser(conn net.Conn) *User {
	u := &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		Conn: conn,
		Msg:  make(chan string),
	}
	// 监听消息
	go u.ListenerMsg()
	return u
}

func (u *User) ListenerMsg() {
	for {
		msg := <-u.Msg
		u.Conn.Write([]byte(msg))
	}
}
