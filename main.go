package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "log"
    "context"
    "errors"
    
    "github.com/graphql-go/graphql"
    "github.com/rs/cors"
    jwt "github.com/dgrijalva/jwt-go"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
    Name: "RootQuery",
    Fields: graphql.Fields{
        "user": &graphql.Field{
            Type: userType,
            Description: "Get single user",
            Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                    Type: graphql.Int,
                },
            },
            Resolve: func(params graphql.ResolveParams) (interface{}, error) {
                id, ok := params.Args["id"].(int)
                if ok {
                    for _, user := range users {
                        if user.ID == id {
                            return user, nil
                        }
                    }
                }
                
                return User{}, nil
            },
        },
        "topic": &graphql.Field{
            Type: topicType,
            Description: "Get single topic",
            Args: graphql.FieldConfigArgument{
                "id": &graphql.ArgumentConfig{
                    Type: graphql.Int,
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                id, ok := p.Args["id"].(int)
                if ok {
                    for _, topic := range topics {
                        if topic.ID == id {
                            return topic, nil
                        }
                    }
                }
                
                return User{}, nil
            },
        },
        "topics": &graphql.Field{
            Type: graphql.NewList(topicType),
            Description: "Get topics",
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                return topics, nil
            },
        },
    },
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
    Name: "RootMutation",
    Fields: graphql.Fields{
        "createUser": &graphql.Field{
            Type: userType,
            Description: "Create new user",
            Args: graphql.FieldConfigArgument{
                "username": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
                "password": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
                "group": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                id := lastUserID + 1
                username, _ := p.Args["username"].(string)
                password, _ := p.Args["password"].(string)
                group, _ := p.Args["group"].(string)
                claims := p.Context.Value("jwtClaims").(jwt.MapClaims)
                
                if group == "admin" {
                    if ok, err := checkPermissions(map[string]bool{
                        "admin": true,
                        "user": false,
                    }, claims, "group"); !ok {
                        return nil, err
                    }
                }
                
                newUser := User{
                    ID: id,
                    Username: username,
                    Password: password,
                    Group: group,
                }
                
                users = append(users, newUser)
                
                lastUserID = id
                
                return newUser, nil
            },
        },
        "createTopic": &graphql.Field{
            Type: topicType,
            Description: "Create new topic",
            Args: graphql.FieldConfigArgument{
                "title": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
                "text": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
                "authorId": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Int),
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                id := lastTopicID + 1
                title, _ := p.Args["title"].(string)
                text, _ := p.Args["text"].(string)
                authorId, _ := p.Args["authorId"].(int)
                claims := p.Context.Value("jwtClaims").(jwt.MapClaims)
                
                if ok, err := checkPermissions(map[string]bool{
                    "admin": true,
                    "user": false,
                }, claims, "group"); !ok {
                    return nil, err
                }
                
                newTopic := Topic{
                    ID: id,
                    Title: title,
                    Text: text,
                    AuthorID: authorId,
                }
                
                topics = append(topics, newTopic)
                
                lastTopicID = id
                
                return newTopic, nil
            },
        },
        "createPost": &graphql.Field{
            Type: postType,
            Description: "Create new post",
            Args: graphql.FieldConfigArgument{
                "text": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.String),
                },
                "authorId": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Int),
                },
                "topicId": &graphql.ArgumentConfig{
                    Type: graphql.NewNonNull(graphql.Int),
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                id := lastPostID + 1
                text, _ := p.Args["text"].(string)
                authorId, _ := p.Args["authorId"].(int)
                topicId, _ := p.Args["topicId"].(int)
                claims := p.Context.Value("jwtClaims").(jwt.MapClaims)
                if len(claims) == 0 {
                    return nil, errors.New("Unauthorized")
                }
                
                if ok, err := checkPermissions(map[string]bool{
                    "admin": true,
                    "user": true,
                }, claims, "group"); !ok {
                    return nil, err
                }
                
                newPost := Post{
                    ID: id,
                    Text: text,
                    AuthorID: authorId,
                    TopicID: topicId,
                }
                
                posts = append(posts, newPost)
                
                lastPostID = id
                
                return newPost, nil
            },
        },
    },
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: rootQuery,
    Mutation: rootMutation,
})

type GraphQLRequest struct {
    Query string `json:"query"`
    Variables map[string]interface{} `json:"variables"`
}

func executeQuery(request GraphQLRequest, schema graphql.Schema, claims jwt.Claims) *graphql.Result {
    result := graphql.Do(graphql.Params{
        Schema: schema,
        RequestString: request.Query,
        VariableValues: request.Variables,
        Context: context.WithValue(context.Background(), "jwtClaims", claims),
    })
    
    if len(result.Errors) > 0 {
        log.Printf("Errors: %+v", result.Errors)
    }
    
    return result
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data, _ := ioutil.ReadAll(r.Body)
        var request GraphQLRequest
        json.Unmarshal(data, &request)
        jwtStr := r.Header.Get("Authorization")
        var (
            claims jwt.MapClaims
            err error
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
        
        res := executeQuery(request, schema, claims)
        
        rjson, _ := json.Marshal(res)
        w.Write(rjson)
    })
    
    mux.HandleFunc("/auth", AuthHandler)
    mux.HandleFunc("/test", TestHandler)
    
    handler := cors.New(cors.Options{
        AllowedHeaders: []string{"Authorization"},
    }).Handler(mux)
    
    http.ListenAndServe(":9000", handler)
}