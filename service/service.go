package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"localProxy/hostsFile"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Service struct {
	httpServerRunning bool
	tlsServerRunning  bool
	// Host-to-host map for default targets
	targetMap map[string]*url.URL

	httpServer *http.Server
	tlsServer  *http.Server

	proxy *httputil.ReverseProxy
}

func (s *Service) StartServer() bool {
	isPrivileged, err := hostsFile.IsPrivileged()
	if err != nil {
		emitErrorToFrontend(err.Error())
		return false
	}

	if !isPrivileged {
		emitWarningToFrontend("You are not running with root/administrator permissions, `apply to /etc/hosts` is skipped!")
	} else {
		// TODO change the /etc/hosts file
	}

	s.makeTargetMap()

	s.proxy = &httputil.ReverseProxy{
		Rewrite: nil,
		Director: func(req *http.Request) {
			hostName := strings.ToLower(req.Host)
			log.Println("Host name is " + hostName)
			target, ok := s.targetMap[hostName]
			if ok {
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
			}
		},
		ModifyResponse: nil,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(502)
			html := fmt.Sprintf("<div style=\"text-align: center\"><h2>Bad Gateway</h2><p>%s</p>", err.Error())
			_, _ = w.Write([]byte(html))
		},
	}

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

	isPrivileged, _ := hostsFile.IsPrivileged()
	if isPrivileged {
		// TODO change the /etc/hosts file
	}
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

func (s *Service) makeTargetMap() {
	s.targetMap = make(map[string]*url.URL)
	for _, h := range config.Hosts {
		if h.DefaultTarget != "" {
			u, err := url.Parse(h.DefaultTarget)
			if err == nil {
				s.targetMap[h.Name] = u
			}
		}
	}
}

func (s *Service) startHttpServer(wg *sync.WaitGroup) {
	s.httpServer = &http.Server{
		Addr:           ":80",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.proxy,
	}

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		emitErrorToFrontend(err.Error())
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
	s.tlsServer = &http.Server{
		Addr:           ":443",
		TLSConfig:      tlsConfig,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        s.proxy,
	}
	listener, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		emitErrorToFrontend(err.Error())
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
