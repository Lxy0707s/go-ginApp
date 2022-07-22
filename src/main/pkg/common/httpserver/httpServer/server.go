package httpServer

import (
	"context"
	"errors"
	"go-ginApp/src/main/pkg/utils/logtool"
	"net/http"
	"time"
)

type (
	// ServerOption is option for Server
	ServerOption struct {
		Address  string
		Handler  http.Handler
		Timeout  int
		CertFile string
		KeyFile  string
	}

	// Server is a HTTP Server
	Server struct {
		log      logtool.Logger
		server   *http.Server
		certFile string
		keyFile  string
	}
)

var (
	// ErrorServerAddrNil is a error which means address is nil
	ErrorServerAddrNil = errors.New("server : addr is nil")
)

// NewServer return Server with option
func NewServer(option ServerOption) *Server {
	svr := &http.Server{
		Addr:              option.Address,
		Handler:           option.Handler,
		ReadHeaderTimeout: time.Second * time.Duration(option.Timeout),
		WriteTimeout:      time.Second * time.Duration(option.Timeout),
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
		IdleTimeout:       time.Minute * 5,
	}
	s := &Server{
		log:      logtool.NewSugar("http-server", false),
		server:   svr,
		certFile: option.CertFile,
		keyFile:  option.KeyFile,
	}
	return s
}

// StartWithHandle is start http server with handleHTTPFailFn
func (s *Server) StartWithHandle(handleHTTPFailFn func()) error {
	if handleHTTPFailFn == nil {
		return errors.New("server.http : handleHTTPFailFn is nil")
	}
	if s.server.Addr == "" {
		return ErrorServerAddrNil
	}
	go func() {
		if s.certFile != "" && s.keyFile != "" {
			s.log.Info("server.http : start https", "addr", s.server.Addr)
			err := s.server.ListenAndServeTLS(s.certFile, s.keyFile)
			if err != http.ErrServerClosed {
				s.log.Error("server.http : start https fail", "addr", s.server.Addr, "error", err)
				handleHTTPFailFn()
			}
		}
		// server.server.ListenAndServeTLS()
		s.log.Info("server.http : start", "addr", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			s.log.Error("server.http : start fail", "addr", s.server.Addr, "error", err)
			handleHTTPFailFn()
		}
	}()
	return nil
}

// Start is start the service
func (s *Server) Start() error {
	if s.server.Addr == "" {
		return ErrorServerAddrNil
	}
	go func() {
		if s.certFile != "" && s.keyFile != "" {
			s.log.Info("server.http : start https", "addr", s.server.Addr)
			err := s.server.ListenAndServeTLS(s.certFile, s.keyFile)
			if err != http.ErrServerClosed {
				s.log.Fatal("server.http : start https fail", "addr", s.server.Addr, "error", err)
			}
		}
		// server.server.ListenAndServeTLS()
		s.log.Info("server.http : start", "addr", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			s.log.Fatal("server.http : start fail", "addr", s.server.Addr, "error", err)
		}
	}()
	return nil
}

// Stop is stop the service
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	s.server.Shutdown(ctx)
	cancel()
	s.server.Close()
	s.log.Info("server.http : stop", "addr", s.server.Addr)
}

type (
	// MultiServerOption is option for MultiServer
	MultiServerOption struct {
		Address   []string
		CertFiles []string
		KeyFiles  []string
		Handler   http.Handler
		Timeout   int
		// TLSConfig *tls.Config
	}

	// MultiServer is a group of HTTP Server
	MultiServer struct {
		log       logtool.Logger
		servers   map[string]*http.Server
		certFiles map[string]string
		keyFiles  map[string]string
	}
)

// NewMultiServer return MultiServer with option
func NewMultiServer(option MultiServerOption) *MultiServer {
	svr := &MultiServer{
		log:       logtool.NewSugar("multi-server", false),
		servers:   make(map[string]*http.Server),
		certFiles: make(map[string]string),
		keyFiles:  make(map[string]string),
	}
	for i, addr := range option.Address {
		svr.certFiles[addr] = option.CertFiles[i]
		svr.keyFiles[addr] = option.KeyFiles[i]
		svr.servers[addr] = &http.Server{
			Addr:              addr,
			Handler:           option.Handler,
			ReadHeaderTimeout: time.Second * time.Duration(option.Timeout),
			WriteTimeout:      time.Second * time.Duration(option.Timeout),
			MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
			IdleTimeout:       time.Minute * 5,
		}
	}
	return svr
}

// Start is start the service
func (ms *MultiServer) Start() error {
	for _, s := range ms.servers {
		if s.Addr == "" {
			return ErrorServerAddrNil
		}
		go func(s *http.Server) {
			certFile, keyFile := ms.certFiles[s.Addr], ms.keyFiles[s.Addr]
			if certFile != "" && keyFile != "" {
				ms.log.Info("server.http : start https", "addr", s.Addr)
				if err := s.ListenAndServeTLS(certFile, keyFile); err != nil {
					if err != http.ErrServerClosed {
						ms.log.Warn("server.http : start https fail", "addr", s.Addr, "error", err)
					}
				}
				return
			}
			ms.log.Info("server.http : start", "addr", s.Addr)
			if err := s.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					ms.log.Warn("server.http : start fail", "addr", s.Addr, "error", err)
				}
			}
		}(s)
	}
	return nil
}

// Stop is stop the service
func (ms *MultiServer) Stop() {
	for _, s := range ms.servers {
		s.Shutdown(context.Background())
		ms.log.Info("server.http : stop", "addr", s.Addr)
	}
}
