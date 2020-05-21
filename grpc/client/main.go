package main

import (
	"context"
	"log"
	"os"
	"time"

	hello "github.com/VicRen/go-playground/grpc"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := hello.NewGreeterClient(conn)
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &hello.HelloRequest{Greeting: name})
	if err != nil {
		panic(err)
	}
	log.Printf("Greeting: %s", r.GetReply())

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err = c.SayHelloAgain(ctx, &hello.HelloRequest{Greeting: name})
	if err != nil {
		panic(err)
	}
	log.Printf("Greeting1: %s", r.GetReply())

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ls, err := c.LotsOfReplies(ctx, &hello.HelloRequest{Greeting: name})
	if err != nil {
		panic(err)
	}
	for {
		r, err := ls.Recv()
		if err != nil {
			break
		}
		log.Printf("Greeting2: %s", r.GetReply())
	}

	gc, err := c.LotsOfGreetings(context.Background())
	if err != nil {
		panic(err)
	}
	count := 5
	for {
		if count < 1 {
			break
		}
		err := gc.Send(&hello.HelloRequest{Greeting: name})
		if err != nil {
			panic(err)
		}
		log.Printf("Sent Greeting3: %d", count)
		count--
	}
	r, err = gc.CloseAndRecv()
	if err != nil {
		panic(err)
	}
	log.Printf("Greeting3: %s", r.GetReply())

	bc, err := c.BidiHello(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		count = 5
		for {
			if count < 1 {
				break
			}
			err = bc.Send(&hello.HelloRequest{Greeting: name})
			if err != nil {
				panic(err)
			}
			count--
		}
		bc.CloseSend()
	}()
	for {
		r, err := bc.Recv()
		if err != nil {
			log.Printf("Greeting4 err: %v", err)
			break
		}
		log.Printf("Greeting4: %s", r.GetReply())
	}
}
