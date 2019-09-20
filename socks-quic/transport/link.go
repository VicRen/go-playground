package transport

import "io"

type Link struct {
	io.Reader
	io.Writer
}
