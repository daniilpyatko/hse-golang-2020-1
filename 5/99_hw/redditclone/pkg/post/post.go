package post

import (
	"redditclone/pkg/comment"
	"redditclone/pkg/user"
	"time"
)

type Vote struct {
	UserId string `json:"user"`
	Vote   int    `json:"vote"`
}

type Post struct {
	Score            int               `json:"score"`
	Views            int               `json:"views"`
	Type             string            `json:"type"`
	Title            string            `json:"title"`
	Author           *user.User        `json:"author"`
	Category         string            `json:"category"`
	Text             string            `json:"text"`
	Url              string            `json:"url"`
	Votes            []Vote            `json:"votes"`
	Comments         []comment.Comment `json:"comments"`
	Created          time.Time         `json:"created"`
	UpvotePercentage int               `json:"upvotePercentage"`
	Id               string            `json:"id"`
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
