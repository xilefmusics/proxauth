package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gobuffalo/packr/v2"

	"proxauth/config"
	"proxauth/login"
	"proxauth/rule"
)

var Config *config.Config
var Html *packr.Box

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	r.URL.Host = r.Host
	rule := rule.Match(Config.Rules, r.URL.Host, r.URL.Path)
	if rule == nil {
		log.Printf("ERROR: No rule found for host=\"%s\" path=\"%s\"", r.URL.Host, r.URL.Path)
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}
	r.URL.Scheme = rule.FromScheme

	if rule.LoginRequired() && rule.IsLoginPath(r.URL.Path) && r.Method == "POST" {
		HandleLoginPOST(w, r)
		return
	}

	if rule.LoginRequired() && rule.IsLoginPath(r.URL.Path) && r.Method == "GET" {
		HandleLoginGET(w, r)
		return
	}

	if rule.LoginRequired() && rule.IsLogoutPath(r.URL.Path) && r.Method == "GET" {
		HandleLogoutGET(w, r, rule)
		return
	}

	if rule.LoginRequired() && HandleCheckLogin(w, r, rule) {
		return
	}

	Forward(w, r, rule)

}

func Forward(w http.ResponseWriter, r *http.Request, rule *rule.Rule) {
	oldUrlString := r.URL.String()
	rule.RewriteRequest(r)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	resp.Body.Close()
	log.Printf("INFO: Forwarded %s to %s", oldUrlString, r.URL.String())
}

func HandleCheckLogin(w http.ResponseWriter, r *http.Request, rule *rule.Rule) bool {
	cookie, err := r.Cookie("proxauth-jwt-token")
	if err != nil {
		if rule.RedirectToLogin {
			http.Redirect(w, r, fmt.Sprintf("%s?redirectedfrom=%s\n", rule.GenLoginUrl(r.URL).String(), r.URL.String()), http.StatusSeeOther)
			log.Printf("INFO: Redirected %s to %s", r.URL.String(), rule.GenLoginUrl(r.URL).String())
			return true
		}
		http.Error(w, "ERROR: Authentication failed!", http.StatusUnauthorized)
		log.Printf("ERROR: Authentication failed! (not exactly one token given)")
		return true
	}

	username, err := login.VerifyJWT(cookie.Value, Config.ServerSecret)
	if err != nil {
		if rule.RedirectToLogin {
			http.Redirect(w, r, fmt.Sprintf("%s?redirectedfrom=%s\n", rule.GenLoginUrl(r.URL).String(), r.URL.String()), http.StatusSeeOther)
			log.Printf("INFO: Redirected %s to %s", r.URL.String(), rule.GenLoginUrl(r.URL).String())
			return true
		}
		http.Error(w, "ERROR: Authentication failed!", http.StatusUnauthorized)
		log.Printf("ERROR: Authentication failed! (%s)", err.Error())
		return true
	}

	if !rule.IsUserPermitted(username) {
		http.Error(w, "ERROR: Authentication failed!", http.StatusUnauthorized)
		log.Printf("ERROR: Authentication failed! (user is not permitted)")
		return true
	}

	r.Header.Del("X-Remote-User")
	r.Header.Set("X-Remote-User", username)

	return false
}

func HandleLoginPOST(w http.ResponseWriter, r *http.Request) {
	username, err := login.Authenticate(Config.Users, r)
	if err != nil {
		log.Printf("ERROR: Login Failed! (%s)", err.Error())
		http.Error(w, "Login failed!", http.StatusUnauthorized)
		return
	}

	expiration := time.Now().UTC().Add(Config.JWTExpirationDuration)

	token, err := login.CreateJWT(username, Config.ServerSecret, expiration.Unix())
	if err != nil {
		log.Printf("ERROR: Login Failed! (%s)", err.Error())
		http.Error(w, "Login Failed!", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(token)
	if err != nil {
		log.Printf("ERROR: Login Failed! (%s)", err.Error())
		http.Error(w, "Login Failed!", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "proxauth-jwt-token", Value: token, Expires: expiration}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write(json)
	log.Printf("INFO: Handled login of user %s for %s", username, r.URL.String())
}

func HandleLogoutGET(w http.ResponseWriter, r *http.Request, rule *rule.Rule) {
	cookie := http.Cookie{Name: "proxauth-jwt-token", Value: "", Expires: time.Unix(0, 0)}
	http.SetCookie(w, &cookie)
	log.Printf("INFO: Handled logout of user %s for %s", "?", r.URL.String())
	http.Redirect(w, r, rule.GenLoginUrl(r.URL).String(), http.StatusSeeOther)
	log.Printf("INFO: Redirected %s to %s", r.URL.String(), rule.GenLoginUrl(r.URL).String())
}

func HandleLoginGET(w http.ResponseWriter, r *http.Request) {
	s, _ := Html.FindString("login.html")
	w.Write([]byte(s))
	log.Printf("INFO: Send back login page for %s", r.URL.String())
}

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatalf("ERROR: Loading config (%s)", err)
	}
	Config = c

	h := packr.New("html", "./html")
	Html = h

	http.HandleFunc("/", HandleRequest)
	log.Fatalf("ERROR: Listening (%s)", http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil))
}
