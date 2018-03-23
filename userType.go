package main

import (
	"context"
	"github.com/graph-gophers/graphql-go"
	"strconv"
)

type user struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"-"`
	Group    string `json:"group"`
}

type userInput struct {
	Username string
	Password string
}

type userResolver struct {
	u user
}

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(r.u.ID))
}

func (r *userResolver) Username() string {
	return r.u.Username
}

func (r *userResolver) Password(ctx context.Context) string {
	return "<hidden>"
}

func (r *userResolver) Group() string {
	return r.u.Group
}
