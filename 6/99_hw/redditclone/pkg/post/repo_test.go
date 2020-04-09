package post

import (
	context "context"
	"errors"
	"redditclone/pkg/comment"
	"redditclone/pkg/random"
	"redditclone/pkg/user"
	"sync"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	mock "github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

// mockgen -source=repo.go -destination=repo_mock.go -package=post CollectionHelper

func TestNewRepo(t *testing.T) {
	NewRepo(&mongo.Collection{})
}

func TestGetAll(t *testing.T) {
	mockCollection := &MockCollection{}
	// var mockCursor CursorHelper
	mockCursor := &MockCursor{}
	// mockSingleResult := &MockSingleResult{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockCollection.
		On("Find", context.TODO(), bson.M{}, options.Find()).
		Return(mockCursor, nil).Once()
	mockCursor.
		On("Next", context.TODO()).
		Return(true).Once()
	mockCursor.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockCursor.
		On("Next", context.TODO()).
		Return(false).Once()
	mockCursor.
		On("Close", context.TODO()).
		Return().Once()
	mockPostRepo.GetAll()
}

func TestGetByCategory(t *testing.T) {
	mockCollection := &MockCollection{}
	mockCursor := &MockCursor{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	category := "programming"
	mockCollection.
		On("Find", context.TODO(), bson.M{"category": category}).
		Return(mockCursor, nil).Once()
	mockCursor.
		On("Next", context.TODO()).
		Return(true).Once()
	mockCursor.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockCursor.
		On("Next", context.TODO()).
		Return(false).Once()
	mockCursor.
		On("Close", context.TODO()).
		Return().Once()
	mockPostRepo.GetByCategory(category)
}

func TestGetById(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockSingleResult := &MockSingleResult{}
	id := "123"

	// OK
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": id}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(nil).Once()

	mockPostRepo.GetById(id)

	// User doesn't exist
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": id}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(errors.New("not found")).Once()
	mockPostRepo.GetById(id)
}

func TestGetByUsername(t *testing.T) {
	mockCollection := &MockCollection{}
	mockCursor := &MockCursor{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	username := "user"
	mockCollection.
		On("Find", context.TODO(), bson.M{"author.username": username}).
		Return(mockCursor, nil).Once()
	mockCursor.
		On("Next", context.TODO()).
		Return(true).Once()
	mockCursor.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockCursor.
		On("Next", context.TODO()).
		Return(false).Once()
	mockCursor.
		On("Close", context.TODO()).
		Return().Once()
	mockPostRepo.GetByUsername(username)
}

func TestAddNewPost(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	curUser := &user.User{}
	curPost := &Post{}
	mockPostRepo.AddNewPost(curUser, curPost)
}

func TestAddPost(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.AddPost(&Post{})
}

func TestAddCommentToPost(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	id := "1"
	filter := bson.M{"id": id}
	mockCollection.
		On("UpdateOne", context.TODO(), filter, mock.Anything).
		Return(nil).Once()
	mockPostRepo.AddCommentToPost(id, &user.User{}, comment.Comment{})
}

func TestPostUpvote(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockSingleResult := &MockSingleResult{}
	curUserId := "1"
	curPost := &Post{}
	curPostId := "23"
	// First vote case

	// calling GetById
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": curPostId}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(func(c interface{}) error {
			return nil
		}).
		Once()

	// calling PostDelete
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": curPostId}).
		Return(nil).Once()

	// calling Addpost
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.PostUpvote(curPostId, curUserId)

	curPost.Votes = append(curPost.Votes, Vote{
		UserId: curUserId,
		Vote:   -1,
	})
	// Repeated vote case

	// calling GetById
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": curPostId}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(func(c interface{}) error {
			cPost := c.(*Post)
			cPost.Votes = append(cPost.Votes, Vote{
				UserId: curUserId,
				Vote:   -1,
			})
			return nil
		}).
		Once()

	// calling PostDelete
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": curPostId}).
		Return(nil).Once()

	// calling Addpost
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.PostUpvote(curPostId, curUserId)
}

func TestPostDownvote(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockSingleResult := &MockSingleResult{}
	curUserId := "1"
	curPost := &Post{}
	curPostId := "23"
	// First vote case

	// calling GetById
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": curPostId}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(func(c interface{}) error {
			return nil
		}).
		Once()

	// calling PostDelete
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": curPostId}).
		Return(nil).Once()

	// calling Addpost
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.PostDownvote(curPostId, curUserId)

	curPost.Votes = append(curPost.Votes, Vote{
		UserId: curUserId,
		Vote:   -1,
	})
	// Repeated vote case

	// calling GetById
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": curPostId}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(func(c interface{}) error {
			cPost := c.(*Post)
			cPost.Votes = append(cPost.Votes, Vote{
				UserId: curUserId,
				Vote:   1,
			})
			return nil
		}).
		Once()

	// calling PostDelete
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": curPostId}).
		Return(nil).Once()

	// calling Addpost
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.PostDownvote(curPostId, curUserId)
}

func TestPostUnvote(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockSingleResult := &MockSingleResult{}
	curUserId := "1"
	curPost := &Post{}
	curPostId := "23"
	// Upvote case

	// calling GetById
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": curPostId}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(func(c interface{}) error {
			cPost := c.(*Post)
			cPost.Votes = append(cPost.Votes, Vote{
				UserId: curUserId,
				Vote:   1,
			})
			return nil
		}).
		Once()

	// calling PostDelete
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": curPostId}).
		Return(nil).Once()

	// calling Addpost
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.PostUnvote(curPostId, curUserId)

	curPost.Votes = append(curPost.Votes, Vote{
		UserId: curUserId,
		Vote:   -1,
	})
	// Downvote case

	// calling GetById
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": curPostId}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(func(c interface{}) error {
			cPost := c.(*Post)
			cPost.Votes = append(cPost.Votes, Vote{
				UserId: curUserId,
				Vote:   -1,
			})
			return nil
		}).
		Once()

	// calling PostDelete
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": curPostId}).
		Return(nil).Once()

	// calling Addpost
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.PostUnvote(curPostId, curUserId)
}

func TestPostDelete(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	id := "1"
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": id}).
		Return(nil).Once()
	mockPostRepo.PostDelete(id)
}

func TestCommentDelete(t *testing.T) {
	mockCollection := &MockCollection{}
	mockPostRepo := PostRepo{
		Collection: mockCollection,
		mu:         &sync.RWMutex{},
		rn:         random.NewGenerator(true),
	}
	mockSingleResult := &MockSingleResult{}
	curUserId := "1"
	curPostId := "23"
	curCommentId := "45"

	// calling GetById
	mockCollection.
		On("FindOne", context.TODO(), bson.M{"id": curPostId}).
		Return(mockSingleResult).Once()
	mockSingleResult.
		On("Decode", mock.AnythingOfType("*post.Post")).
		Return(func(c interface{}) error {
			cPost := c.(*Post)
			cPost.Comments = append(cPost.Comments, comment.Comment{
				Author: &user.User{
					Id: curUserId,
				},
				Id: curCommentId,
			})
			return nil
		}).
		Once()

	// calling PostDelete
	mockCollection.
		On("DeleteOne", context.TODO(), bson.M{"id": curPostId}).
		Return(nil).Once()

	// calling Addpost
	mockCollection.
		On("InsertOne", context.TODO(), mock.AnythingOfType("*post.Post")).
		Return(nil).Once()
	mockPostRepo.CommentDelete(curPostId, curCommentId, curUserId)
}

// go test -v  -coverprofile=repo.out && go tool cover -html=repo.out -o repo.html && rm repo.out
