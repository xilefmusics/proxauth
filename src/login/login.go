package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Salt     string `json:"salt" yaml:"salt"`
}

func ValidataUser(users []User, checkUser User) bool {
	var user *User = nil
	for _, u := range users {
		if u.Username == checkUser.Username {
			user = &u
		}
	}
	if user == nil {
		log.Println("user not found")
		return false
	}

	// TODO: Check password

	return true
}

func Authenticate(users []User, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		return "", errors.New("http method is not of kind post")
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return "", err
	}

	if user.Username == "" {
		return "", errors.New("no username given")
	}

	if user.Password == "" {
		return "", errors.New("no password given")
	}

	if !ValidataUser(users, user) {
		return "", errors.New("user and password don't match")
	}

	return user.Username, nil
}

func CreateJWT(username string, serverSecret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(10 * time.Hour).Unix()
	claims["user"] = username
	return token.SignedString(serverSecret)
}

func VerifyJWT(tokenString string, serverSecret []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return serverSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%s", claims["user"]), nil
}
