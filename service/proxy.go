package service

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyCore struct {
	proxy    *httputil.ReverseProxy
	rulesMap map[string][]*HostProxy
}

func NewProxyCore() *ProxyCore {
	core := new(ProxyCore)
	core.proxy = &httputil.ReverseProxy{
		Rewrite:        nil,
		Director:       core.proxyDirector,
		ModifyResponse: nil,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(502)
			html := fmt.Sprintf("<div style=\"text-align: center\"><h2>Bad Gateway</h2><p>%s</p>", err.Error())
			_, _ = w.Write([]byte(html))
		},
	}

	return core
}

func (p *ProxyCore) SetRulesMap(rules map[string][]*HostProxy) {
	p.rulesMap = rules
}

func (p *ProxyCore) proxyDirector(req *http.Request) {
	hostName := strings.ToLower(req.Host)
	log.Println("Host name is " + hostName)
	rules, ok := p.rulesMap[hostName]
	if !ok {
		return
	}

rulesLoop:
	for _, rule := range rules {
		if rule.targetUrl == nil {
			u, err := url.Parse(rule.Target)
			if err != nil {
				// Ignore current rule when target is invalid
				continue
			}
			rule.targetUrl = u
		}

		if p.isRuleMatched(req.URL, rule) {
			req.URL.Scheme = rule.targetUrl.Scheme
			req.URL.Host = rule.targetUrl.Host

			req.Header.Add("X-Forwarded-For", req.RemoteAddr)
			req.Header.Add("X-Forwarded-Proto", req.URL.Scheme)
			req.Header.Add("X-Forwarded-Host", req.Host)

			if !rule.ChangeOrigin {
				req.Host = rule.targetUrl.Host
			}

			log.Println(req.URL.String())

			break rulesLoop
		}

	}
}

func (p *ProxyCore) isRuleMatched(url *url.URL, rule *HostProxy) bool {
	switch rule.MatchType {
	case MatchTypePrefix:
		return strings.HasPrefix(url.Path, rule.MatchRule)

	default:
		return false
	}
}

func (p *ProxyCore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
