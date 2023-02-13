package server

import (
	"fmt"
	"net"
	"strings"
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
	msg = strings.TrimSpace(msg)
	fmt.Println(msg)
	if msg == "who" {
		replayMsg := ""
		for name, _ := range u.Server.OnlineMap {
			if u.Name != name {
				userList := fmt.Sprintf("[%s] is online... \n", name)
				replayMsg += userList
			}
		}
		SendMsg(replayMsg, u)
	} else if len(msg) >= 7 && string(msg[:7]) == "rename|" {
		// 1.用户名校验并查询名称是否存在
		name := msg[7:]
		if len(name) == 0 {
			SendMsg("username is not empty\n", u)
			return
		}
		_, ok := u.Server.OnlineMap[name]
		if ok {
			SendMsg("username is exist\n", u)
			return
		}
		// 2.更新用户名
		u.Server.Lock.Lock()
		delete(u.Server.OnlineMap, u.Name)
		u.Name = name
		u.Server.OnlineMap[u.Name] = u

		// 3.返回更新成功消息
		SendMsg("username is update \n", u)
		u.Server.Lock.Unlock()

	} else if len(msg) >= 3 && string(msg[:3]) == "to|" {
		msgSplit := strings.Split(msg[3:], "|")
		// 简单实现，会有bug（消息中存在 ”|“ ）
		if len(msgSplit) != 2 {
			SendMsg("msg format error, please input fix format,exp: [to|username|msg]\n", u)
			return
		}
		u.Server.Lock.Lock()
		toUserName := msgSplit[0]
		toUser, ok := u.Server.OnlineMap[toUserName]
		if !ok {
			SendMsg("username is not exist\n", u)
			u.Server.Lock.Unlock()
			return
		}
		SendMsg(fmt.Sprintf("%s say: %s \n", u.Name, msgSplit[1]), toUser)
		u.Server.Lock.Unlock()
	} else {
		u.Server.BroadCast(u, msg)
	}

}
func (u *User) ListenerMsg() {
	for {
		msg, ok := <-u.Msg
		if !ok {
			fmt.Println("user chan is close!")
		} else {
			u.Conn.Write([]byte(msg))
		}
	}
}

// 这里可以对消息返回一个布尔值或者返回一个err，来表示消息发送是否成功
func SendMsg(msg string, user *User) {
	_, err := user.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("消息发送失败")
	}
}
