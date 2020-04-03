package post

import (
	"redditclone/pkg/comment"
	"sync"
)

type PostRepo struct {
	data []*Post
	mu   *sync.RWMutex
}

func NewRepo() *PostRepo {
	return &PostRepo{
		data: make([]*Post, 0, 10),
		mu:   &sync.RWMutex{},
	}
}

func (p *PostRepo) GetAll() ([]*Post, error) {
	return p.data, nil
}

func (p *PostRepo) GetByCategory(category string) ([]*Post, error) {
	res := make([]*Post, 0, 10)
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < len(p.data); i++ {
		if p.data[i].Category == category {
			res = append(res, p.data[i])
		}
	}
	return res, nil
}

func (p *PostRepo) GetById(id string) (*Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < len(p.data); i++ {
		if p.data[i].Id == id {
			return p.data[i], nil
		}
	}
	return nil, nil
}

func (p *PostRepo) GetByUsername(username string) ([]*Post, error) {
	res := make([]*Post, 0, 10)
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < len(p.data); i++ {
		if p.data[i].Author.Username == username {
			res = append(res, p.data[i])
		}
	}
	return res, nil
}

func (p *PostRepo) AddPost(post *Post) {
	p.mu.Lock()
	p.data = append(p.data, post)
	p.mu.Unlock()
}

func (p *PostRepo) AddCommentToPost(postId string, comment comment.Comment) {
	post, _ := p.GetById(postId)
	p.mu.Lock()
	post.Comments = append(post.Comments, comment)
	p.mu.Unlock()
}

func (p *PostRepo) DeletePost(postId string) {
	postInd := -1
	p.mu.Lock()
	for ind, val := range p.data {
		if val.Id == postId {
			postInd = ind
		}
	}
	p.mu.Unlock()
	if postInd != -1 {
		p.mu.Lock()
		p.data = append(p.data[:postInd], p.data[postInd+1:]...)
		p.mu.Unlock()
	}
}

func (p *PostRepo) PostUpvote(postId, curUserId string) {
	post, _ := p.GetById(postId)
	found := false
	p.mu.Lock()
	for i := 0; i < len(post.Votes); i++ {
		if post.Votes[i].UserId == curUserId {
			if post.Votes[i].Vote == -1 {
				post.Votes[i].Vote = 1
				post.Score += 2
				post.UpvotePercentage = post.GetUpvotePercentage()
				found = true
			}
		}
	}
	if !found {
		newVote := Vote{
			UserId: curUserId,
			Vote:   1,
		}
		post.Votes = append(post.Votes, newVote)
		post.Score += 1
		post.UpvotePercentage = post.GetUpvotePercentage()
	}
	p.mu.Unlock()
}

func (p *PostRepo) PostDownvote(postId, curUserId string) {
	post, _ := p.GetById(postId)
	found := false
	p.mu.Lock()
	for i := 0; i < len(post.Votes); i++ {
		if post.Votes[i].UserId == curUserId {
			if post.Votes[i].Vote == 1 {
				post.Votes[i].Vote = -1
				post.Score -= 2
				post.UpvotePercentage = post.GetUpvotePercentage()
				found = true
			}
		}
	}
	if !found {
		newVote := Vote{
			UserId: curUserId,
			Vote:   -1,
		}
		post.Votes = append(post.Votes, newVote)
		post.Score -= 1
		post.UpvotePercentage = post.GetUpvotePercentage()
	}
	p.mu.Unlock()
}

func (p *PostRepo) PostUnvote(postId, curUserId string) {
	post, _ := p.GetById(postId)
	eraseInd := -1
	p.mu.Lock()
	for i := 0; i < len(post.Votes); i++ {
		if post.Votes[i].UserId == curUserId {
			if post.Votes[i].Vote == 1 {
				post.Score -= 1
			} else {
				post.Score += 1
			}
			eraseInd = i
		}
	}
	if eraseInd != -1 {
		post.Votes = append(post.Votes[:eraseInd], post.Votes[eraseInd+1:]...)
		post.UpvotePercentage = post.GetUpvotePercentage()
	}
	p.mu.Unlock()
}

func (p *PostRepo) CommentDelete(postId, commentId, curUserId string) {
	post, _ := p.GetById(postId)
	commentInd := -1
	p.mu.Lock()
	for i := 0; i < len(post.Comments); i++ {
		if post.Comments[i].Id == commentId && post.Comments[i].Author.Id == curUserId {
			commentInd = i
		}
	}
	if commentInd != -1 {
		post.Comments = append(post.Comments[:commentInd], post.Comments[commentInd+1:]...)
	}
	p.mu.Unlock()
}

func (p *PostRepo) PostDelete(postId string) {
	postInd := -1
	p.mu.Lock()
	for i := 0; i < len(p.data); i++ {
		if p.data[i].Id == postId {
			postInd = i
		}
	}
	if postInd != -1 {
		p.data = append(p.data[:postInd], p.data[postInd+1:]...)
	}
	p.mu.Unlock()
}
