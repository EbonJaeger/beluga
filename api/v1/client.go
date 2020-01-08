package v1

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Client is a client for the API
type Client struct {
	client *http.Client
}

// NewClient will return a new client for our Unix socket for
// communication with the daemon
func NewClient(address string) *Client {
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", address)
				},
				DisableKeepAlives:     false,
				IdleConnTimeout:       30 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 60 * time.Second,
		},
	}
}

// Close will kill any idle connections to avoid leaking file
// descriptors
func (c *Client) Close() {
	c.client.Transport.(*http.Transport).CloseIdleConnections()
}
