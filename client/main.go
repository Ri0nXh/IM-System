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
	client.Run()

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
		Name:       fmt.Sprintf("%s:%d", ip, port),
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

func (c *Client) Menu() {

	menu := `================= Menu =================
1 => Common chat
2 => Private chat
3 => Rename name
4 => Show Menu
0 => Quit program
========================================
`
	fmt.Println(menu)
}

func (c *Client) ReviceUserInput() int {
	var action int
	_, err := fmt.Scanln(&action)
	if err != nil {
		fmt.Println("input number is error")
		return -1
	}
	if action >= 0 && action <= 3 {
		return action
	} else {
		return -1
	}

}

func (c *Client) Run() {
	c.Menu()
	// 接收用户输入
	for {
		// 接收用户输入
		action := c.ReviceUserInput()
		fmt.Printf("[%s] choice %d \n", c.Name, action)
		switch action {
		case 0:
		case 1:
		case 2:
		case 3:
		case 4:
			c.Menu()
		case -1:
			c.Menu()
		}
	}
}
