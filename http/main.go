package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintln(writer, "<h1>Hello world from Vic</h1>")
	})
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
