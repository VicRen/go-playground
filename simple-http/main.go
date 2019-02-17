package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("%s %s\nHost: %s\n\n", request.Method, request.Proto, request.Host)
		io.WriteString(writer, fmt.Sprintf("Request %s\n", request.Method))
	})
	http.ListenAndServe(":8081", nil)
}
