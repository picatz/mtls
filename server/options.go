package server

import "crypto/tls"

// Options contains each available configuration option
// for an mTLS SSH Server.
type Options struct {
	Addr      string
	TLSConfig *tls.Config
	Handler   func(*tls.Conn)
}

// Option implements a hook to custom a Server
// using the New function.
type Option func(*Options) error

// WithAddr sets the given addr string for the Server.
func WithAddr(addr string) Option {
	return func(o *Options) error {
		o.Addr = addr
		return nil
	}
}

// WithTLSConfig sets the given TLS config for the Server.
func WithTLSConfig(config *tls.Config) Option {
	return func(o *Options) error {
		o.TLSConfig = config
		return nil
	}
}

// WithHandler sets the server's handler function.
func WithHandler(h func(*tls.Conn)) Option {
	return func(o *Options) error {
		o.Handler = h
		return nil
	}
}
