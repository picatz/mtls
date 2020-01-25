package server

import (
	"crypto/tls"
	"log"
	"net"
	"time"
)

// Server implements an mTLS SSH server.
type Server struct {
	started   bool
	addr      string
	tlsConfig *tls.Config
	listener  net.Listener
	handler   func(*tls.Conn)
	shutdown  chan struct{}
}

// New implements a wrappeer to create a new Server,
// applying the given server Option(s).
func New(opts ...Option) (*Server, error) {
	serverOptions := &Options{
		Addr: DefaultAddr,
	}

	server := &Server{}

	for _, opt := range opts {
		err := opt(serverOptions)
		if err != nil {
			return nil, err
		}
	}

	server.addr = serverOptions.Addr
	server.tlsConfig = serverOptions.TLSConfig
	server.handler = serverOptions.Handler

	server.shutdown = make(chan struct{})

	listener, err := tls.Listen("tcp", server.addr, server.tlsConfig)
	if err != nil {
		return nil, err
	}
	server.listener = listener

	return server, nil
}

// Listener returns the underlying net.Listener from the Server.
func (s *Server) Listener() net.Listener {
	return s.listener
}

// HandleConn will handle a connection from the server's accept loop.
func (s *Server) HandleConn(conn *tls.Conn) {
	if s.handler != nil {
		s.handler(conn)
	}
}

// Start will start the server's accept loop.
func (s *Server) Start() {
	go func() {
		log.Println("server: started")
		for {
			select {
			case <-s.shutdown:
				log.Println("server: stopping")
				return
			default:
				conn, err := s.listener.Accept()
				if err != nil {
					log.Printf("failed to accept conn: %s", err)
					time.Sleep(2 * time.Second)
					break
				}
				log.Println("server: accepted connection")
				tlsConn, ok := conn.(*tls.Conn)
				if !ok {
					log.Fatalf("failed to cast conn to tls.Conn")
				}

				log.Println("server: handling tls conn")
				go s.HandleConn(tlsConn)
			}
		}
	}()
}

// Shutdown stops the underlying accept loop.
func (s *Server) Shutdown() net.Listener {
	return s.listener
}
