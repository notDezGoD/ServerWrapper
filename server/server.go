package socketserver

import (
	"fmt"
	"strconv"

	"golang.org/x/sys/unix"
)

type Server struct {
	addr     *unix.SockaddrInet4
	socket   int
	recieve  chan string
	send     chan string
	shutdown bool
	nfd      int
}

func (s *Server) Init() {
	s.recieve = make(chan string, 10)
	s.send = make(chan string)
	s.shutdown = false
	s.nfd = 0
}
func (s *Server) CreateSocket() error {
	var err error
	s.socket, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	return err
}

func (s *Server) Close() {
	unix.Close(s.socket)
}

func (s *Server) CreateAddress(addr [4]byte, port int) {
	s.addr = &unix.SockaddrInet4{Port: port, Addr: addr}
}

func (s *Server) Connect() error {
	err := unix.Connect(s.socket, s.addr)
	return err
}

func (s *Server) Listen() {
	buff := make([]byte, 1024)

	for !s.shutdown {
		err := unix.Listen(s.socket, 10)
		if err != nil {
			panic(err)
		}

		fmt.Println("Listening on port :" + strconv.Itoa(s.addr.Port))

		nfd, _, err := unix.Accept(s.socket) // not getting addr
		if err != nil {
			panic(err)
		}

		n, err := unix.Read(nfd, buff)
		if err != nil {
			panic(err)
		}
		s.recieve <- string(buff[:n])
	}
	fmt.Println("Stopped Listening")
}

func (s *Server) Broadcastmsg(msg string) {
	s.send <- msg
}

func (s *Server) Recievemsg() {
	for !s.shutdown {
		for {
			value, ok := <-s.recieve
			if !ok {
				fmt.Println("Error on recieving from channel")
			} else {
				fmt.Println(value)
			}
		}
	}
}
func (s *Server) Sendmsg() {
	for !s.shutdown {
		if s.nfd > 0 {
			for {
				value, ok := <-s.send
				if !ok {
					fmt.Println("Error on geting send value")
				} else {
					res, _ := unix.Write(s.nfd, []byte(value))
					fmt.Println("Send Result :" + strconv.Itoa(res))
				}
			}
		}
	}
}
func (s *Server) Bind() error {
	return unix.Bind(s.socket, s.addr)
}
