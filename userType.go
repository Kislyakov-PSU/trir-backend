package main

import (
    "github.com/graphql-go/graphql"
)

type User struct {
    ID int `json:"id"`
    Username string `json:"username"`
    Password string
    Group string `json:"group"`
}

var userType = graphql.NewObject(graphql.ObjectConfig{
    Name: "User",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "username": &graphql.Field{
            Type: graphql.String,
        },
        "group": &graphql.Field{
            Type: graphql.String,
        },
    },
})