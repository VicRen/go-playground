package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	url string
)

type HttpClient interface {
	Get(string) (*http.Response, error)
}

func init() {
	flag.StringVar(&url, "url", "http://baidu.com", "Which url do we want to parse?")
	flag.Parse()
}

func send(hc HttpClient, link string) error {
	response, err := hc.Get(link)
	if err != nil {
		return err
	}

	if response == nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

func main() {
	client := &http.Client{}
	err := send(client, url)

	if err != nil {
		fmt.Println(err)
	}
}
