package service

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ProxyCore struct {
	rulesMap map[string][]*HostProxy
}

func NewProxyCore() *ProxyCore {
	core := new(ProxyCore)

	return core
}

func (p *ProxyCore) SetRulesMap(rules map[string][]*HostProxy) {
	for _, p := range rules {
		for _, rule := range p {
			switch rule.TargetType {
			case TargetTypeStatic:
				t := StaticTarget{}
				_ = mapstructure.Decode(rule.TargetParams, &t)
				rule.targetObj = &t

			case TargetTypeProxy:
				t := ProxyTarget{}
				_ = mapstructure.Decode(rule.TargetParams, &t)
				rule.targetObj = &t
			case TargetTypePHP:
				t := PhpTarget{}
				_ = mapstructure.Decode(rule.TargetParams, &t)
				rule.targetObj = &t
			default:
				rule.targetObj = nil
			}
		}
	}
	p.rulesMap = rules
}

func (p *ProxyCore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hostName := strings.ToLower(r.Host)
	if hostName == "" {
		hostName = strings.ToLower(r.URL.Host)
	}
	log.Println("Host name is " + hostName)
	rules, ok := p.rulesMap[hostName]
	if !ok {
		return
	}

rulesLoop:
	for _, rule := range rules {
		if p.isRuleMatched(r.URL, rule) {
			//log.Printf("Rule matched, %s, %s\n", rule.TargetType, utils.JsonEncode(rule))

			if rule.targetObj != nil {
				rule.targetObj.ServeTarget(rule, w, r)
			}

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
