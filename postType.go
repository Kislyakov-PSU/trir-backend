package main

import (
	"github.com/graph-gophers/graphql-go"
	"strconv"
)

type post struct {
	ID      int    `json:"id" gorm:"primary_key"`
	Text    string `json:"text" gorm:"size:5000"`
	TopicID int
	UserID  int
}

type postInput struct {
	Text    string
	TopicID graphql.ID
	UserID  graphql.ID
}

type postResolver struct {
	p post
}

func (r *postResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(r.p.ID))
}

func (r *postResolver) Text() string {
	return r.p.Text
}

func (r *postResolver) User() *userResolver {
	var _user user
	db.Model(&r.p).Related(&_user)
	return &userResolver{_user}
}
