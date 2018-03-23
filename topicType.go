package main

import (
	"github.com/graph-gophers/graphql-go"
	"strconv"
)

type topic struct {
	ID     int    `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Text   string `json:"text" gorm:"size:5000"`
	UserID int
}

type topicInput struct {
	Title  string
	Text   string
	UserID graphql.ID
}

type topicResolver struct {
	t topic
}

func (r *topicResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(r.t.ID))
}

func (r *topicResolver) Title() string {
	return r.t.Title
}

func (r *topicResolver) Text() string {
	return r.t.Text
}

func (r *topicResolver) User() *userResolver {
	var _user user
	db.Model(&r.t).Related(&_user)
	return &userResolver{_user}
}

func (r *topicResolver) Posts() *[]*postResolver {
	var list []*postResolver
	var posts []post

	db.Model(&r.t).Related(&posts)

	for _, post := range posts {
		list = append(list, &postResolver{post})
	}

	return &list
}
