package socketserver

import (
	"fmt"
	"strconv"

	"socketwrapper/manager" // your client manager package

	"golang.org/x/sys/unix"
)

type Server struct {
	addr   *unix.SockaddrInet4
	socket int
	cm     *manager.ClientManager
}

func (s *Server) Init(cm *manager.ClientManager) {
	s.cm = cm
}

func (s *Server) CreateSocket() error {
	var err error
	s.socket, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	err1 := unix.SetsockoptInt(s.socket, unix.IPPROTO_TCP, unix.TCPOPT_MAXSEG, 1300)
	if err1 != nil {
		fmt.Println(err)
	}
	return err
}

func (s *Server) CreateAddress(addr [4]byte, port int) {
	s.addr = &unix.SockaddrInet4{Port: port, Addr: addr}
}

func (s *Server) Bind() error {
	return unix.Bind(s.socket, s.addr)
}

func (s *Server) Listen() error {
	if err := unix.Listen(s.socket, 10); err != nil {
		return err
	}
	fmt.Println("Listening on port :" + strconv.Itoa(s.addr.Port))
	return nil
}

func (s *Server) AcceptLoop() {
	for {
		fd, _, err := unix.Accept(s.socket)
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		fmt.Println("New client connected, fd:", fd)

		// Add to manager
		s.cm.Add(fd)

		// Optionally start goroutine to read client messages
		go s.handleClient(fd)
	}
}

func (s *Server) handleClient(fd int) {
	buff := make([]byte, 1300)
	for {
		n, err := unix.Read(fd, buff)
		if err != nil {
			fmt.Println("Client disconnected, fd:", fd)
			s.cm.Remove(fd)
			return
		}
		msg := string(buff[:n])
		if msg != "" {
			fmt.Printf("Received from %d: %s\n", fd, msg)
		}

		// Example: broadcast to everyone
		s.cm.Broadcast([]byte(msg))
	}
}
