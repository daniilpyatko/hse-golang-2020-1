package post

import (
	"context"
	"errors"
	"redditclone/pkg/comment"
	"redditclone/pkg/random"
	"redditclone/pkg/user"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepo struct {
	Collection CollectionHelper
	mu         *sync.RWMutex
	rn         *random.Generator
}

func NewRepo(ccol *mongo.Collection) *PostRepo {
	return &PostRepo{
		Collection: &mongoCollection{coll: ccol},
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
}

func (p *PostRepo) GetAll() ([]*Post, error) {
	var result []*Post
	cur, _ := p.Collection.Find(context.TODO(), bson.M{}, options.Find())
	for cur.Next(context.TODO()) {
		var elem Post
		cur.Decode(&elem)
		result = append(result, &elem)
	}
	cur.Close(context.TODO())
	return result, nil
}

func (p *PostRepo) GetByCategory(category string) ([]*Post, error) {
	var result []*Post
	cur, _ := p.Collection.Find(context.TODO(), bson.M{"category": category})
	for cur.Next(context.TODO()) {
		var elem Post
		cur.Decode(&elem)
		result = append(result, &elem)
	}
	cur.Close(context.TODO())
	return result, nil
}

func (p *PostRepo) GetById(id string) (*Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var result Post
	err := p.Collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	if err != nil {
		return nil, errors.New("User with given Id doesn't exist")
	}
	return &result, nil
}

func (p *PostRepo) GetByUsername(username string) ([]*Post, error) {
	var result []*Post
	cur, _ := p.Collection.Find(context.TODO(), bson.M{"author.username": username})
	for cur.Next(context.TODO()) {
		var elem Post
		cur.Decode(&elem)
		result = append(result, &elem)
	}
	cur.Close(context.TODO())
	return result, nil
}

func (p *PostRepo) AddNewPost(user *user.User, newPost *Post) {
	newPost.Author = user
	newPost.Votes = []Vote{
		Vote{
			UserId: newPost.Author.Id,
			Vote:   1,
		},
	}
	newPost.Score = 1
	newPost.Views = 0
	newPost.Created = time.Now()
	newPost.UpvotePercentage = 100
	newPost.Comments = make([]comment.Comment, 0, 10)
	newPost.Id = p.rn.GetString()
	p.Collection.InsertOne(context.TODO(), newPost)
}

func (p *PostRepo) AddPost(post *Post) {
	p.Collection.InsertOne(context.TODO(), post)
}

func (p *PostRepo) AddCommentToPost(postId string, curUser *user.User, curComment comment.Comment) {
	curComment = comment.Comment{
		Created: time.Now(),
		Author:  curUser,
		Body:    curComment.Body,
		Id:      p.rn.GetString(),
	}
	filter := bson.M{"id": postId}
	update := bson.M{"$push": bson.M{"comments": curComment}}
	p.Collection.UpdateOne(context.TODO(), filter, update)
}

func (p *PostRepo) PostUpvote(postId, curUserId string) {
	post, _ := p.GetById(postId)
	found := false
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
	p.PostDelete(postId)
	p.AddPost(post)
}

func (p *PostRepo) PostDownvote(postId, curUserId string) {
	post, _ := p.GetById(postId)
	found := false
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
	p.PostDelete(postId)
	p.AddPost(post)
}

func (p *PostRepo) PostUnvote(postId, curUserId string) {
	post, _ := p.GetById(postId)
	eraseInd := -1
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
	p.PostDelete(postId)
	p.AddPost(post)
}

func (p *PostRepo) CommentDelete(postId, commentId, curUserId string) {
	post, _ := p.GetById(postId)
	for i, curComment := range post.Comments {
		if curComment.Author.Id == curUserId && curComment.Id == commentId {
			post.Comments = append(post.Comments[:i], post.Comments[i+1:]...)
			break
		}
	}
	p.PostDelete(postId)
	p.AddPost(post)
}

func (p *PostRepo) PostDelete(postId string) {
	filter := bson.M{"id": postId}
	p.Collection.DeleteOne(context.TODO(), filter)
}
