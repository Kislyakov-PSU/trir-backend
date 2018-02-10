package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "log"
    
    "github.com/graphql-go/graphql"
    "github.com/rs/cors"
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

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query: rootQuery,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
    result := graphql.Do(graphql.Params{
        Schema: schema,
        RequestString: query,
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
        
        res := executeQuery(string(data), schema)
        
        rjson, _ := json.Marshal(res)
        w.Write(rjson)
    })
    
    mux.HandleFunc("/auth", AuthHandler)
    mux.HandleFunc("/test", TestHandler)
    
    handler := cors.Default().Handler(mux)
    
    http.ListenAndServe(":9000", handler)
}