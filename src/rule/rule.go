package rule

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Rule struct {
	FromScheme      string   `json:"fromScheme" yaml:"fromScheme"`
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

func (self *Rule) SetDefaults() {
	if self.FromScheme == "" {
		self.FromScheme = "http"
	}
	if self.FromHost == "" {
		self.FromHost = "*"
	}
	if len(self.FromPath) > 0 && self.FromPath[len(self.FromPath)-1] == '/' {
		self.FromPath = self.FromPath[:len(self.FromPath)-1]
	}
	if self.ToScheme == "" {
		self.ToScheme = "http"
	}
	if self.ToHost == "" {
		self.ToHost = "localhost"
	}
	if len(self.ToPath) > 0 && self.ToPath[len(self.ToPath)-1] == '/' {
		self.ToPath = self.ToPath[:len(self.ToPath)-1]
	}
	if self.LoginPath == "" {
		self.LoginPath = "/login"
	}
	if self.LogoutPath == "" {
		self.LogoutPath = "/logout"
	}
	if self.AllowedUsers == nil {
		self.AllowedUsers = []string{}
	}
}

func (self *Rule) IsLoginPath(path string) bool {
	return path == self.FromPath+self.LoginPath
}

func (self *Rule) IsLogoutPath(path string) bool {
	return path == self.FromPath+self.LogoutPath
}

func (self *Rule) Match(fromHost, fromPath string) bool {
	if len(fromPath) == 0 || fromPath[len(fromPath)-1] != '/' {
		fromPath += "/"
	}
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
	return len(self.AllowedUsers) > 0
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

func (self *Rule) GenLoginUrl(url *url.URL) *url.URL {
	newUrl := *url
	newUrl.Path = self.FromPath + self.LoginPath
	return &newUrl
}

func Match(rules []Rule, fromHost, fromPath string) *Rule {
	for _, rule := range rules {
		if rule.Match(fromHost, fromPath) {
			return &rule
		}
	}
	return nil
}
