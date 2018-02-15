package main

import (
    "errors"
    jwt "github.com/dgrijalva/jwt-go"
)

func jwtExtractClaims(tokenStr string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenStr, func (token *jwt.Token) (interface{}, error) {
        return pubKey, nil
    })
    
    if token == nil {
        return jwt.MapClaims{}, errors.New("Empty jwt token")
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    } else {
        return jwt.MapClaims{}, err
    }
}

func checkPermissions(permissions map[string]bool, claims jwt.MapClaims, key string) (bool, error) {
    if len(claims) == 0 {
        return false, errors.New("Unauthorized")
    } else if permissions[claims[key].(string)] {
        return true, nil
    } else {
        return false, errors.New("Unprivileged")
    }
}