package service

import (
	"fmt"
	"github.com/yookoala/gofast"
	"localProxy/utils"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type TargetHandler interface {
	ServeTarget(rule *HostProxy, w http.ResponseWriter, r *http.Request)
}

type StaticTarget struct {
	DocumentRoot string `json:"documentRoot"`
	Fallback     string `json:"fallback"`
}

type PhpTarget struct {
	DocumentRoot string `json:"documentRoot"`
	FPMAddress   string `json:"fpmAddress"`
	Fallback     string `json:"fallback"`

	phpHandler      gofast.Handler
	fallbackHandler gofast.Handler
}

type ProxyTarget struct {
	ProxyPass    string `json:"proxyPass"`
	ChangeOrigin bool   `json:"changeOrigin"`
	proxyPassUrl *url.URL
	proxy        *httputil.ReverseProxy
}

func (s *StaticTarget) ServeTarget(rule *HostProxy, w http.ResponseWriter, r *http.Request) {
	staticFilePath := tryStaticFile(s.DocumentRoot, r.URL.Path)
	if staticFilePath != "" {
		log.Println(staticFilePath)
		http.ServeFile(w, r, staticFilePath)
		return
	}

	w.WriteHeader(404)
	_, _ = w.Write([]byte("File not found"))
}

func (t *ProxyTarget) ServeTarget(rule *HostProxy, w http.ResponseWriter, r *http.Request) {
	if t.proxy == nil {
		t.proxy = &httputil.ReverseProxy{
			Rewrite: nil,
			Director: func(req *http.Request) {
				if t.proxyPassUrl == nil {
					u, err := url.Parse(t.ProxyPass)
					if err != nil {
						return
					}
					t.proxyPassUrl = u
				}

				req.URL.Scheme = t.proxyPassUrl.Scheme
				req.URL.Host = t.proxyPassUrl.Host

				req.Header.Add("X-Forwarded-For", req.RemoteAddr)
				req.Header.Add("X-Forwarded-Proto", req.URL.Scheme)
				req.Header.Add("X-Forwarded-Host", req.Host)

				if t.ChangeOrigin {
					req.Host = t.proxyPassUrl.Host
				}
			},
			ModifyResponse: nil,
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				w.WriteHeader(502)
				html := fmt.Sprintf("<div style=\"text-align: center\"><h2>Bad Gateway</h2><p>%s</p>", err.Error())
				_, _ = w.Write([]byte(html))
			},
		}
	}

	t.proxy.ServeHTTP(w, r)
}

func (p *PhpTarget) ServeTarget(rule *HostProxy, w http.ResponseWriter, r *http.Request) {
	// init
	if p.phpHandler == nil {
		if p.FPMAddress == "" {
			p.FPMAddress = "127.0.0.1:9000"
		}
		connFactory := gofast.SimpleConnFactory("tcp", p.FPMAddress)
		sess := gofast.Chain(gofast.BasicParamsMap, gofast.MapHeader, gofast.MapRemoteHost)(gofast.BasicSession)

		p.phpHandler = gofast.NewHandler(gofast.NewPHPFS(p.DocumentRoot)(sess), gofast.SimpleClientFactory(connFactory))
		if p.Fallback != "" {
			p.fallbackHandler = gofast.NewHandler(
				gofast.NewFileEndpoint(singleJoiningSlash(p.DocumentRoot, p.Fallback))(sess),
				gofast.SimpleClientFactory(connFactory),
			)
		}
	}

	// PHP file handler
	if strings.HasSuffix(strings.ToLower(r.URL.Path), ".php") {
		p.phpHandler.ServeHTTP(w, r)
		return
	}

	// Static file handler
	staticFilePath := tryStaticFile(p.DocumentRoot, r.URL.Path)
	if staticFilePath != "" {
		log.Println(staticFilePath)
		http.ServeFile(w, r, staticFilePath)
		return
	}

	// Fallback
	if p.Fallback != "" {
		p.fallbackHandler.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(404)
	_, _ = w.Write([]byte("File not found"))
}

func tryStaticFile(docRoot, filePath string) string {
	if strings.HasSuffix(filePath, "/") {
		filePath = filePath[0 : len(filePath)-1]
	}
	path := docRoot + filePath
	if utils.IsDir(path) {
		path = path + "/index.html"
		if utils.FileExists(path) {
			return path
		}
		return ""
	}

	if utils.FileExists(path) {
		return path
	}
	return ""
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
