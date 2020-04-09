package comment

import (
	"redditclone/pkg/user"
	"time"
)

type Comment struct {
	Created time.Time  `json:"created" bson:"created"`
	Author  *user.User `json:"author" bson:"author"`
	Body    string     `json:"body" bson:"body"`
	Id      string     `json:"id" bson:"id"`
}
