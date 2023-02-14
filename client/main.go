package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	serverIp   string
	serverPort int
)

// 初始化控制台输入参数信息
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "please input your ip address")
	flag.IntVar(&serverPort, "port", 8888, "please input your port")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	err := client.Connect()
	if err != nil {
		fmt.Println("connect server error: ", err)
	} else {
		fmt.Println("connect server success")
	}

	// 启动客户端业务

}

type Client struct {
	ServerIp   string
	Name       string
	ServerPort int
	Conn       net.Conn
}

// 初始化客户端
func NewClient(ip string, port int) *Client {
	c := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}
	return c
}

// 建立连接对象
func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ServerIp, c.ServerPort))
	if err != nil {
		fmt.Println("conn server error: ", err)
	}
	c.Conn = conn
	return err
}
