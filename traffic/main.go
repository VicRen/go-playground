package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	command "github.com/VicRen/go-playground/traffic/pb"
	"google.golang.org/grpc"
)

var (
	target = flag.String("t", "", "")
)

func main() {
	flag.Parse()

	if *target == "" {
		panic("invalid target")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(fmt.Sprintf("failed to dial: %s: %v", *target, err))
	}
	defer conn.Close()

	qs := command.NewStatsServiceClient(conn)

	req := &command.QueryStatsRequest{
		Pattern: "user>>>",
		Reset_:  true,
	}

	resp, err := qs.QueryStats(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println("----------response-------")
	for _, r := range resp.Stat {
		fmt.Println("------>Stat: Name:", r.Name, "Value:", r.Value)
	}

}
