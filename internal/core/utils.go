package core

import (
	"fmt"
	"github.com/chainreactors/zombie/pkg"
	"net/url"
	"strings"
)

type Target struct {
	IP      string            `json:"ip"`
	Port    string            `json:"port"`
	Service pkg.Service       `json:"service"`
	Param   map[string]string `json:"param"`
}

func (t *Target) String() string {
	return fmt.Sprintf("%s://%s:%s", t.Service, t.IP, t.Port)
}

func (t *Target) UpdateService(s string) {
	t.Service = pkg.Service(strings.ToLower(s))
	if t.Port == "" {
		t.Port = t.Service.DefaultPort()
	}
}

func (t *Target) Addr() *utils.Addr {
	return &utils.Addr{IP: utils.NewIP(t.IP), Port: t.Port}
}

func ParseUrl(u string) (*Target, bool) {
	parsed, err := url.Parse(u)
	if err != nil {
		return nil, false
	}
	if parsed.Host == "" {
		if utils.IsIpv4(u) {
			return &Target{
				IP: u,
			}, true
		}
		return nil, false
	}
	t := &Target{
		IP: parsed.Hostname(),
	}
	if parsed.Port() != "" {
		t.Port = parsed.Port()
	}

	if parsed.Scheme != "" {
		t.Service = pkg.Service(parsed.Scheme)
		if t.Port == "" {
			t.Port = t.Service.DefaultPort()
		}
	} else if t.Port != "" {
		t.Service = pkg.GetDefault(t.Port)
	}

	return t, true
}
