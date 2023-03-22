package rule

import (
	"fmt"
	"net/http"
	"strings"
)

type Rule struct {
	// TODO set defaults
	// TODO Differentiate between API and User page
	// TODO URL tools
	FromHost        string   `json:"fromHost" yaml:"fromHost"`
	FromPath        string   `json:"fromPath" yaml:"fromPath"`
	ToScheme        string   `json:"toScheme" yaml:"toScheme"`
	ToHost          string   `json:"toHost" yaml:"toHost"`
	ToPort          int      `json:"toPort" yaml:"toPort"`
	ToPath          string   `json:"toPath" yaml:"toPath"`
	LoginPath       string   `json:"loginPath" yaml:"loginPath"`
	LogoutPath      string   `json:"logoutPath" yaml:"logoutPath"`
	AllowedUsers    []string `json:"allowedUsers" yaml:"allowedUsers"`
	RedirectToLogin bool     `json:"redirectToLogin" yaml:"redirectToLogin"`
}

func (self *Rule) IsLoginPath(path string) bool {
	return self.RewritePath(path) == self.LoginPath
}

func (self *Rule) IsLogoutPath(path string) bool {
	return self.RewritePath(path) == self.LogoutPath
}

func (self *Rule) Match(fromHost, fromPath string) bool {
	return (self.FromHost == fromHost || self.FromHost == "*") && strings.HasPrefix(fromPath, self.FromPath)
}

func (self *Rule) IsUserPermitted(username string) bool {
	for _, user := range self.AllowedUsers {
		if user == username {
			return true
		}
	}
	return false
}

func (self *Rule) LoginRequired() bool {
	return self.LoginPath != ""
}

func (self *Rule) RewritePath(path string) string {
	return strings.Replace(path, self.FromPath, self.ToPath, 1)
}

func (self *Rule) RewriteRequest(r *http.Request) {
	r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	r.URL.Scheme = self.ToScheme
	r.URL.Host = fmt.Sprintf("%s:%d", self.ToHost, self.ToPort)
	r.URL.Path = self.RewritePath(r.URL.Path)
	r.RequestURI = ""
}

func (self *Rule) GenLoginUrl(host string) string {
	return fmt.Sprintf("%s%s", host, self.LoginPath)
}

func Match(rules []Rule, fromHost, fromPath string) *Rule {
	for _, rule := range rules {
		if rule.Match(fromHost, fromPath) {
			return &rule
		}
	}
	return nil
}
