package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	hello "github.com/VicRen/go-play-ground/grpc"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	log.Printf("listening on: %s", l.Addr())

	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, &server{})
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}

type server struct {
}

func (s *server) BidiHello(bs hello.Greeter_BidiHelloServer) error {
	count := 0
	for {
		r, err := bs.Recv()
		if err != nil {
			log.Printf("BidiHello Recived err: %v", err)
			return err
		}
		count++
		resp := r.GetGreeting()
		bs.Send(&hello.HelloResponse{Reply: resp})
		log.Printf("BidiHello Recived: %v, count: %d", r.GetGreeting(), count)
	}
	return nil
}

func (s *server) LotsOfGreetings(gs hello.Greeter_LotsOfGreetingsServer) error {
	var resp string
	count := 0
	for {
		r, err := gs.Recv()
		if err != nil {
			log.Printf("LotsOfGreetings Recived err: %v", err)
			gs.SendAndClose(&hello.HelloResponse{Reply: "Hello " + resp})
			return nil
		}
		count++
		log.Printf("LotsOfGreetings Recived: %v, count: %d", r.GetGreeting(), count)
		resp += " " + r.GetGreeting()
	}
}

func (s *server) LotsOfReplies(in *hello.HelloRequest, rs hello.Greeter_LotsOfRepliesServer) error {
	log.Printf("LotsOfReplies Recived: %v", in.GetGreeting())
	count := 5
	for count > 0 {
		err := rs.Send(&hello.HelloResponse{Reply: fmt.Sprintf("Hello %s:%d", in.GetGreeting(), count)})
		if err != nil {
			return err
		}
		count--
	}
	return io.EOF
}

func (s *server) SayHelloAgain(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	log.Printf("SayHelloAgain Recived: %v", in.GetGreeting())
	return &hello.HelloResponse{Reply: "Hello again " + in.GetGreeting()}, nil
}

func (s *server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	log.Printf("SayHello Recived: %v", in.GetGreeting())
	return &hello.HelloResponse{Reply: "Hello " + in.GetGreeting()}, nil
}
