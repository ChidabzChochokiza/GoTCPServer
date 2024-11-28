package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	wg         sync.WaitGroup
	quit       chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quit:       make(chan struct{}),
	}
}
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return nil
	}

	s.listener = listener
	fmt.Println("Server started on ")

	for {
		select {
		case <-s.quit:
			fmt.Println("Stoping server...")
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return nil
				default:
					fmt.Println("Accepting error: ", err)
					continue
				}
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer s.wg.Done()

	buf := make([]byte, 2024)

	for {
		select {
		case <-s.quit:
			return
		default:
			readBytes, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Client disconnected")
				} else {
					fmt.Println("Error reading from connection")
				}
				return
			}

			msg := buf[:readBytes]
			fmt.Println(string(msg))
		}
	}
}
func (s *Server) Stop() {
	close(s.quit)

	if s.listener != nil {
		s.listener.Close()
	}

	s.wg.Wait()
	fmt.Println("Server stopped")
}

func main() {
	server := NewServer(":8080")

	go func() {
		if err := server.Start(); err != nil {
			fmt.Println("Error starting server", err)
		}
	}()

	fmt.Println("Press Enter to stop server")
	fmt.Scanln()

	server.Stop()
}
