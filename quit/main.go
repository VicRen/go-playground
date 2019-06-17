package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Service struct {
	ch chan struct{}
	waitGroup *sync.WaitGroup
}

func NewService() *Service {
	s := &Service{
		make(chan struct{}),
		&sync.WaitGroup{},
	}
	s.waitGroup.Add(1)
	return s
}

func (s *Service) Serve(listener *net.TCPListener) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("stopping listening on", listener.Addr())
			err := listener.Close()
			if err != nil {
				fmt.Println(err)
			}
			return
		default:
		}
		err := listener.SetDeadline(time.Now().Add(1e9))
		if err != nil {
			fmt.Println(err)
		}
		conn, err := listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			fmt.Println(err)
		}
		log.Println(conn.RemoteAddr(), "connected")
		s.waitGroup.Add(1)
		go s.serve(conn)
	}
}

func (s *Service) Stop() {
	close(s.ch)
	s.waitGroup.Wait()
}

func (s *Service) serve(conn *net.TCPConn) {
	defer conn.Close()
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}
		conn.SetDeadline(time.Now().Add(1e9))
		buf := make([]byte, 4096)
		if _, err := conn.Read(buf); err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
			return
		}
		if _, err := conn.Write(buf); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:48879")
	if err != nil {
		log.Println(err)
	}
	listener, err := net.ListenTCP("tcp", laddr)

	service:=NewService()
	go service.Serve(listener)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	service.Stop()
}
