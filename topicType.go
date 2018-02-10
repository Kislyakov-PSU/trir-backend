package main

import (
    "github.com/graphql-go/graphql"
)

type Topic struct {
    ID int `json:"id"`
    Title string `json:"title"`
    Text string `json:"text"`
    AuthorID int
}

var topicType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Topic",
    Fields: graphql.Fields{
        "id": &graphql.Field{
            Type: graphql.Int,
        },
        "title": &graphql.Field{
            Type: graphql.String,
        },
        "text": &graphql.Field{
            Type: graphql.String,
        },
        "author": &graphql.Field{
            Type: userType,
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                if topic, ok := p.Source.(Topic); ok {
                    for _, user := range users {
                        if user.ID == topic.AuthorID {
                            return user, nil
                        }
                    }
                }
                return User{}, nil
            },
        },
        "posts": &graphql.Field{
            Type: graphql.NewList(postType),
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                if topic, ok := p.Source.(Topic); ok {
                    var _posts []Post
                    for _, post := range posts {
                        if post.TopicID == topic.ID {
                            _posts = append(_posts, post)
                        }
                    }
                    return _posts, nil
                }
                return []interface{}{}, nil
            },
        },
    },
})