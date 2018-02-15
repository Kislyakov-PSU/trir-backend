package main

var users = []User{
    User{ID: 1, Username: "Defman", Group: "admin", Password: "password"},
    User{ID: 2, Username: "Dafmen", Group: "user", Password: "possward"},
}

var posts = []Post{
    Post{ID: 1, TopicID: 1, Text: "hello", AuthorID: 1},
    Post{ID: 2, TopicID: 1, Text: "world", AuthorID: 2},
    Post{ID: 3, TopicID: 1, Text: "world1", AuthorID: 2},
    Post{ID: 4, TopicID: 2, Text: "world2", AuthorID: 2},
    Post{ID: 5, TopicID: 2, Text: "world3", AuthorID: 1},
    Post{ID: 6, TopicID: 2, Text: "world4", AuthorID: 1},
}

var topics = []Topic{
    Topic{ID: 1, Title: "Hello world!", Text: "A simple topic", AuthorID: 1},
    Topic{ID: 2, Title: "World hello!", Text: "Not so simple ehh?", AuthorID: 3},
}

var (
    lastUserID int = 2
    lastPostID int = 6
    lastTopicID int = 2
)