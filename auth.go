package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type authResponse struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
	Error  string `json:"error"`
	User   user   `json:"user"`
}

const (
	privKeyPath = "key.pem"
	pubKeyPath  = "jwt.pem"
)

var (
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
)

func init() {
	privBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	privKey, err = jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		log.Fatal(err)
	}

	pubBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	pubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		log.Fatal(err)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)

	info := strings.Split(string(data), ";")
	login := info[0]
	password := info[1]

	var _user user
	db.Where("username = ? AND password = ?", login, password).First(&_user)

	if _user.ID == 0 {
		msg, _ := json.Marshal(&authResponse{
			Status: 500,
			Error:  "Invalid credentials",
		})
		w.Write(msg)
		return
	}

	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["group"] = _user.Group
	claims["username"] = _user.Username
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	tokenString, err := token.SignedString(privKey)
	if err != nil {
		msg, _ := json.Marshal(&authResponse{
			Status: 500,
			Error:  fmt.Sprintf("%v", err),
		})
		w.Write(msg)
		log.Print(err)
		return
	}
	msg, _ := json.Marshal(&authResponse{
		Status: 200,
		Token:  tokenString,
		User:   _user,
	})
	w.Write(msg)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	tokenStr := string(data)

	claims, err := jwtExtractClaims(tokenStr)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(fmt.Sprintf("%+v", claims)))
	}
}
