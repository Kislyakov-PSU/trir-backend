package main

import (
    "github.com/graphql-go/graphql"
)

type Post struct {
    ID int `json:"id"`
    Text string `json:"text"`
    TopicID int
    AuthorID int
}

var postType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Post",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "text": &graphql.Field{
            Type: graphql.String,
        },
        "author": &graphql.Field{
            Type: userType,
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                if post, ok := p.Source.(Post); ok {
                    for _, user := range users {
                        if user.ID == post.AuthorID {
                            return user, nil
                        }
                    }
                }
                return User{}, nil
            },
        },
    },
})