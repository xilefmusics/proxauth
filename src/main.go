package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"proxauth/config"
	"proxauth/login"
	"proxauth/rule"
)

var Config *config.Config

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// TODO: proper handling
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rule := rule.Match(Config.Rules, r.URL.Host, r.URL.Path)
	if rule == nil {
		log.Printf("ERROR: No rule found for host=\"%s\" path=\"%s\"", r.URL.Host, r.URL.Path)
		http.Error(w, "404 Not found", http.StatusNotFound)
		return
	}

	log.Printf("INFO: Handle Request method=%s clientHost=%s fromHost=%s fromPath=%s toHost=%s:%d toPath=%s", r.Method, r.RemoteAddr, r.URL.Host, r.URL.Path, rule.ToHost, rule.ToPort, rule.ToPath)

	if rule.LoginRequired() {

		if rule.LoginPath == rule.RewritePath(r.URL.Path) {
			username, err := login.Authenticate(Config.Users, r) // Users is empty
			if err != nil {
				log.Printf("ERROR: Login Failed! (%s)", err.Error())
				http.Error(w, "Login failed!", http.StatusUnauthorized)
				return
			}

			token, err := login.CreateJWT(username, Config.ServerSecret)
			if err != nil {
				log.Printf("ERROR: Login Failed! (%s)", err.Error())
				http.Error(w, "Login Failed!", http.StatusUnauthorized)
				return
			}

			json, err := json.Marshal(token)
			if err != nil {
				log.Printf("ERROR: Login Failed! (%s)", err.Error())
				http.Error(w, "Login Failed!", http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(json)
			return
		}

		if len(r.Header["Token"]) != 1 {
			http.Error(w, "ERROR: Authentication failed!", http.StatusUnauthorized)
			log.Printf("ERROR: Authentication failed! (not exactly one token given)")
			return
		}

		username, err := login.VerifyJWT(r.Header["Token"][0], Config.ServerSecret)
		if err != nil {
			http.Error(w, "ERROR: Authentication failed!", http.StatusUnauthorized)
			log.Printf("ERROR: Authentication failed! (%s)", err.Error())
			return
		}

		if !rule.IsUserPermitted(username) {
			http.Error(w, "ERROR: Authentication failed!", http.StatusUnauthorized)
			log.Printf("ERROR: Authentication failed! (user is not permitted)")
			return
		}

		r.Header.Del("user")
		r.Header.Set("user", username)

	}

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

}

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatalf("ERROR: Loading config (%s)", err)
	}
	Config = c

	http.HandleFunc("/", HandleRequest)
	log.Fatalf("ERROR: Listening (%s)", http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil))
}
