package main

import (
	clientmanager "socketwrapper/manager"
	socketserver "socketwrapper/server"
)

func main() {
	var s socketserver.Server
	cm := clientmanager.NewClientManager()
	s.Init(cm)
	err := s.CreateSocket()
	if err != nil {
		panic(err)
	}
	s.CreateAddress([4]byte{127, 0, 0, 1}, 8080)
	err = s.Bind()
	if err != nil {
		panic(err)
	}
	s.Listen()
	s.AcceptLoop()

	select {}
}
