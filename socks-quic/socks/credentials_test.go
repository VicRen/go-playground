package socks

import (
	"testing"
)

func TestStaticCredentials(t *testing.T) {
	sc := StaticCredentials{
		"foo": "bar",
		"baz": "",
	}

	if !sc.Valid("foo", "bar") {
		t.Fatalf("expect valid")
	}

	if !sc.Valid("baz", "") {
		t.Fatalf("expect valid")
	}

	if sc.Valid("foo", "") {
		t.Fatalf("expect invalid")
	}
}
