package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"net"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	httpServerRunning bool
	tlsServerRunning  bool

	httpServer *http.Server
	tlsServer  *http.Server

	proxyCore *ProxyCore
}

func (s *Service) StartServer() bool {
	go func() {
		if err := addHostsFileRecord(); err != nil {
			emitErrorToFrontend(err.Error())
		}
	}()

	s.proxyCore = NewProxyCore()
	s.proxyCore.SetRulesMap(s.makeRulesMap())

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Start http server
	go s.startHttpServer(wg)

	// Start https server
	go s.startTlsServer(wg)

	wg.Wait()

	return s.httpServerRunning || s.tlsServerRunning
}

func (s *Service) ShutdownServer() {
	if s.httpServer != nil {
		err := s.httpServer.Shutdown(context.Background())
		if err != nil {
			emitErrorToFrontend(err.Error())
		}
		s.httpServerRunning = false
	}
	if s.tlsServer != nil {
		err := s.tlsServer.Shutdown(context.Background())
		if err != nil {
			emitErrorToFrontend(err.Error())
		}
		s.tlsServerRunning = false
	}

	go func() {
		if err := removeHostsFileRecord(); err != nil {
			emitErrorToFrontend(err.Error())
		}
	}()
}

func (s *Service) GetServerStatus() map[string]bool {
	return map[string]bool{
		"HTTP": s.httpServerRunning,
		"TLS":  s.tlsServerRunning,
	}
}

func (s *Service) loadTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	// Load certificate files. Invalid certificate will be ignored and emit a warning
	tlsConfig.Certificates = make([]tls.Certificate, 0, len(config.Hosts))
	for _, h := range config.Hosts {
		if h.EnableTLS && h.TLSCertFile != "" && h.TLSKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(h.TLSCertFile, h.TLSKeyFile)
			if err != nil {
				emitWarningToFrontend(fmt.Sprintf("Load certificate for %s failed\n", h.Name))
			} else {
				tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
			}
		}
	}

	return tlsConfig
}

func (s *Service) makeRulesMap() map[string][]*HostProxy {
	result := make(map[string][]*HostProxy)
	for _, h := range config.Hosts {
		result[h.Name] = h.Proxies
	}
	return result
}

func (s *Service) startHttpServer(wg *sync.WaitGroup) {
	s.httpServer = &http.Server{
		Addr:           ":" + cast.ToString(config.Setting.HttpPort),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.proxyCore,
	}

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		emitErrorToFrontend(err.Error())
		wg.Done()
		return
	}
	s.httpServerRunning = true
	wg.Done()
	err = s.httpServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		emitErrorToFrontend(err.Error())
	}
	s.httpServerRunning = false
}

func (s *Service) startTlsServer(wg *sync.WaitGroup) {
	tlsConfig := s.loadTLSConfig()
	// Skip TLS server when no certificate available
	if len(tlsConfig.Certificates) == 0 {
		s.tlsServerRunning = false
		wg.Done()
		return
	}

	s.tlsServer = &http.Server{
		Addr:           ":" + cast.ToString(config.Setting.HttpsPort),
		TLSConfig:      tlsConfig,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.proxyCore,
	}
	listener, err := tls.Listen("tcp", s.tlsServer.Addr, tlsConfig)
	if err != nil {
		emitErrorToFrontend(err.Error())
		wg.Done()
		return
	}
	s.tlsServerRunning = true
	wg.Done()
	err = s.tlsServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		emitErrorToFrontend(err.Error())
	}
	s.tlsServerRunning = false
}
