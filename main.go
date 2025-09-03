package main

import (
	socketserver "socketwrapper/server"
)

func main() {
	var s socketserver.Server
	defer s.Close()

	s.Init()
	s.CreateAddress([4]byte{127, 0, 0, 1}, 8080)
	s.CreateSocket()
	s.Bind()
	//s.Connect()
	go s.Recievemsg()
	go s.Listen()
	select {}

}
