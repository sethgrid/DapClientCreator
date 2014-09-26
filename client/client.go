package DapClient

import (
	"net/http"
)

type Client struct {
	DapAddr    string
	HttpClient *http.Client
}

func New(dapUrl string) (*Client, error) {
	client := &Client{DapAddr: dapUrl, HttpClient: &http.Client{}}
	return client, nil
}

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type ErrorResponse struct {
	Error string
}
