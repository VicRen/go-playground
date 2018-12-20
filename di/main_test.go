package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockHttpClient struct {
}

func (m *MockHttpClient) Get(string) (*http.Response, error) {
	response := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte("Test Response"))),
	}
	return response, nil
}

func TestSendWithValidResponse(t *testing.T) {
	httpClient := &MockHttpClient{}
	err := send(httpClient, "something")
	if err != nil {
		t.Errorf("Shouldn't have received an error with a valid MockHttpClient, got %s", err)
	}
}
