package client

import (
	"crypto/tls"
)

// Client implements an mTLS SSH client.
type Client struct {
	addr      string
	tlsConfig *tls.Config
}

// New implements a wrappeer to create a new Client,
// applying the given client Option(s).
func New(opts ...Option) (*Client, error) {
	clientOptions := &Options{}

	client := &Client{}

	for _, opt := range opts {
		err := opt(clientOptions)
		if err != nil {
			return nil, err
		}
	}

	client.addr = clientOptions.Addr
	client.tlsConfig = clientOptions.TLSConfig

	client.tlsConfig.BuildNameToCertificate()

	return client, nil
}

func (c *Client) Dial() (*tls.Conn, error) {
	return tls.Dial("tcp", c.addr, c.tlsConfig)
}
