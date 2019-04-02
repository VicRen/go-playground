package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	deadline := 1
	d := time.Now().Add(time.Duration(deadline) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	resChan := make(chan string)

	go HttpDoTest(ctx, resChan)

	var resData string
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.Tick(time.Duration(time.Duration(deadline*2) * time.Second)):
		fmt.Println("time out!")
	case resData = <-resChan:
		fmt.Println("Read data finished")
	}
	log.Printf("Read data size: [%d]", len(resData))
}

func HttpDoTest(ctx context.Context, resChan chan<- string) error {
	start := time.Now()

	repoUrl := "https://baidu.com"
	req, err := http.NewRequest("GET", repoUrl, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest error: %v", err)
	}

	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do error: %v", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll error: %v", err)
	}
	log.Printf("Read body size [%d]", len(data))
	log.Printf("CostTime is: %s", time.Since(start).String())
	resChan <- string(data)
	return nil
}
