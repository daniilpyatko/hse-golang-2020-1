package comment

import (
	"redditclone/pkg/user"
	"time"
)

type Comment struct {
	Created time.Time  `json:"created"`
	Author  *user.User `json:"author"`
	Body    string     `json:"body"`
	Id      string     `json:"id"`
}
