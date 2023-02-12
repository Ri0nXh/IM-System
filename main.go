package main

import "IM-System/server"

//import svc "IM-System/im-server"

func main() {
	s := server.NewServer("127.0.0.1", 8888)
	s.Start()
}
