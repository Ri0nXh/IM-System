package main

import "IM-System/im-server"

//import svc "IM-System/im-server"

func main() {
	s := im_server.NewServer("127.0.0.1", 8888)
	s.Start()
}
