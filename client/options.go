package client

import "crypto/tls"

// Options contains each available configuration option
// for an mTLS SSH Client.
type Options struct {
	Addr      string
	TLSConfig *tls.Config
}

// Option implements a hook to custom a Client
// using the New function.
type Option func(*Options) error

// WithAddr sets the given addr string for the Client.
func WithAddr(addr string) Option {
	return func(o *Options) error {
		o.Addr = addr
		return nil
	}
}

// WithTLSConfig sets the given TLS config for the Client.
func WithTLSConfig(config *tls.Config) Option {
	return func(o *Options) error {
		o.TLSConfig = config
		return nil
	}
}
