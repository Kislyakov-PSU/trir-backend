package main

import (
	"context"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graph-gophers/graphql-go"
	"github.com/rs/cors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type graphQLRequest struct {
	Query         string
	OperationName string
	Variables     map[string]interface{}
}

type key int

const jwtClaimsKey key = 0
const schema = `
schema {
	query: Query
	mutation: Mutation
}

type Query {
	topics: [Topic]
	topic(id: ID!): Topic
	user(id: ID!): User
}

type Mutation {
	createUser(user: UserInput!): User
	createTopic(topic: TopicInput!): Topic
	createPost(post: PostInput!): Post
}

type User {
	id: ID!
	username: String!
	password: String!
	group: String!
}

input UserInput {
	username: String!
	password: String!
}

type Topic {
	id: ID!
	title: String!
	text: String!
	user: User!
	posts: [Post]
}

input TopicInput {
	title: String!
	text: String!
	userID: ID!
}

type Post {
	id: ID!
	text: String!
	user: User!
}

input PostInput {
	userID: ID!
	topicID: ID!
	text: String!
}
`

type resolver struct{}

func (r *resolver) Topics() *[]*topicResolver {
	var topics []topic
	var list []*topicResolver

	db.Find(&topics)

	for _, topic := range topics {
		list = append(list, &topicResolver{topic})
	}
	return &list
}

func (r *resolver) Topic(args struct{ ID graphql.ID }) *topicResolver {
	var _topic topic
	db.Where("id = ?", args.ID).First(&_topic)
	return &topicResolver{_topic}
}

func (r *resolver) User(ctx context.Context, args struct{ ID graphql.ID }) *userResolver {
	var _user user
	db.Where("id = ?", args.ID).First(&_user)
	return &userResolver{_user}
}

func (r *resolver) CreateUser(args *struct {
	User *userInput
}) *userResolver {
	_user := &user{
		Username: args.User.Username,
		Password: args.User.Password,
		Group:    "user",
	}

	db.Create(_user)

	return &userResolver{*_user}
}

func (r *resolver) CreateTopic(ctx context.Context, args *struct {
	Topic *topicInput
}) *topicResolver {
	claims := ctx.Value(jwtClaimsKey).(jwt.MapClaims)
	log.Printf("Claims: %+v", claims)

	if ok, _ := checkPermissions(map[string]bool{
		"admin": true,
		"user":  false,
	}, claims, "group"); !ok {
		return nil
	}

	uid, _ := strconv.Atoi(string(args.Topic.UserID))
	_topic := &topic{
		Title:  args.Topic.Title,
		Text:   args.Topic.Text,
		UserID: uid,
	}

	db.Create(_topic)

	log.Printf("%+v", _topic)

	return &topicResolver{*_topic}
}

func (r *resolver) CreatePost(ctx context.Context, args *struct {
	Post *postInput
}) *postResolver {
	claims := ctx.Value(jwtClaimsKey).(jwt.MapClaims)
	log.Printf("%+v", claims)

	if ok, _ := checkPermissions(map[string]bool{
		"admin": true,
		"user":  true,
	}, claims, "group"); !ok {
		return nil
	}

	uid, _ := strconv.Atoi(string(args.Post.UserID))
	tid, _ := strconv.Atoi(string(args.Post.TopicID))
	_post := &post{
		Text:    args.Post.Text,
		UserID:  uid,
		TopicID: tid,
	}

	db.Create(_post)

	return &postResolver{*_post}
}

func main() {
	connect()
	graphqlSchema := graphql.MustParseSchema(schema, &resolver{})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		var request graphQLRequest
		json.Unmarshal(data, &request)
		jwtStr := r.Header.Get("Authorization")

		var (
			claims jwt.MapClaims
			err    error
		)

		if jwtStr != "" {
			log.Printf("Auth header: %v", jwtStr)
			claims, err = jwtExtractClaims(jwtStr)
			if err != nil {
				log.Printf("Auth error: %v", err)
			}
		} else {
			claims = jwt.MapClaims{}
		}

		ctx := context.WithValue(r.Context(), jwtClaimsKey, claims)

		res := graphqlSchema.Exec(ctx, request.Query, request.OperationName, request.Variables)
		rjson, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		w.Write(rjson)
	})

	mux.HandleFunc("/auth", authHandler)

	handler := cors.New(cors.Options{
		AllowedHeaders: []string{"Authorization"},
	}).Handler(mux)

	http.ListenAndServe(":9000", handler)
}
