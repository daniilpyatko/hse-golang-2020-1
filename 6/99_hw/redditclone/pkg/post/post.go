package post

import (
	"redditclone/pkg/comment"
	"redditclone/pkg/user"
	"time"
)

type Vote struct {
	UserId string `json:"user" bson:"user"`
	Vote   int    `json:"vote" bson:"vote"`
}

type Post struct {
	Score            int               `json:"score" bson:"score"`
	Views            int               `json:"views" bson:"views"`
	Type             string            `json:"type" bson:"type"`
	Title            string            `json:"title"  bson:"title"`
	Author           *user.User        `json:"author" bson:"author"`
	Category         string            `json:"category" bson:"category"`
	Text             string            `json:"text" bson:"text"`
	Url              string            `json:"url" bson:"url"`
	Votes            []Vote            `json:"votes" bson:"votes"`
	Comments         []comment.Comment `json:"comments" bson:"comments"`
	Created          time.Time         `json:"created" bson:"created"`
	UpvotePercentage int               `json:"upvotePercentage" bson:"upvotePercentage"`
	Id               string            `json:"id" bson:"id"`
}

func (p *Post) GetUpvotePercentage() int {
	cntMinus := 0
	cntPlus := 0
	for i := 0; i < len(p.Votes); i++ {
		if p.Votes[i].Vote == -1 {
			cntMinus += 1
		} else {
			cntPlus += 1
		}
	}
	if cntMinus+cntPlus == 0 {
		return 0
	} else {
		return (cntPlus * 100) / (cntMinus + cntPlus)
	}
}
