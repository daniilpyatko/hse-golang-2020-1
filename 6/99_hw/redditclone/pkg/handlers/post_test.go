package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"redditclone/pkg/comment"
	"redditclone/pkg/post"
	"redditclone/pkg/user"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestPostHandlerAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// registered user
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}
	// test := map[string]string{
	// 	"username": "user",
	// 	"password": "pass",
	// 	"userid":   "1",
	// }
	postRepoResult := []*post.Post{}
	postRepoMock.EXPECT().GetAll().Return(postRepoResult, nil)
	// jsonQuery, _ := json.Marshal(query)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	service.AllPosts(w, r)
}

func TestPostHandlerPostsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}
	postRepoResult := []*post.Post{}
	postRepoMock.EXPECT().GetByCategory("programming").Return(postRepoResult, nil)
	r := httptest.NewRequest("GET", "/", nil)
	r = mux.SetURLVars(r, map[string]string{"category": "programming"})
	w := httptest.NewRecorder()
	service.PostsByCategory(w, r)
}

func TestPostHandlerPostsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}
	postRepoResult := &post.Post{}
	postRepoMock.EXPECT().GetById("1").Return(postRepoResult, nil)
	postRepoMock.EXPECT().PostDelete("1").Return()
	postRepoMock.EXPECT().AddPost(postRepoResult).Return()
	r := httptest.NewRequest("GET", "/", nil)
	r = mux.SetURLVars(r, map[string]string{"post_id": "1"})
	w := httptest.NewRecorder()
	service.PostsById(w, r)
}

func TestPostHandlerPostsByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}
	postRepoResult := []*post.Post{}
	postRepoMock.EXPECT().GetByUsername("user").Return(postRepoResult, nil)
	r := httptest.NewRequest("GET", "/", nil)
	r = mux.SetURLVars(r, map[string]string{"username": "user"})
	w := httptest.NewRecorder()
	service.PostsByUsername(w, r)
}

func TestPostHandlerNewPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}

	// user does not exist
	r := httptest.NewRequest("POST", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), "Id", "1"))
	w := httptest.NewRecorder()
	userRepoMock.EXPECT().GetUserById("1").Return(nil, user.ErrUserNotFound)
	service.NewPost(w, r)

	// OK
	user := &user.User{
		Username: "user",
		Id:       "1",
	}
	userRepoMock.EXPECT().GetUserById(user.Id).Return(user, nil)
	query := map[string]string{}
	jsonQuery, _ := json.Marshal(query)
	newPost := post.Post{}

	postRepoMock.EXPECT().AddNewPost(user, &newPost).Return()
	r = httptest.NewRequest("POST", "/", bytes.NewReader(jsonQuery))
	r = r.WithContext(context.WithValue(r.Context(), "Id", "1"))
	w = httptest.NewRecorder()
	service.NewPost(w, r)
}

func TestPostHandlerAddCommentToPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}

	// OK
	query := map[string]string{
		"comment": "This is a comment.",
	}
	curUser := &user.User{
		Id:       "1",
		Username: "user",
	}
	curComment := comment.Comment{
		Body: query["comment"],
	}
	curPost := post.Post{}
	userRepoMock.EXPECT().GetUserById("1").Return(curUser, nil)
	postRepoMock.EXPECT().AddCommentToPost("2", curUser, curComment).Return()
	postRepoMock.EXPECT().GetById("2").Return(&curPost, nil)
	jsonQuery, _ := json.Marshal(query)
	r := httptest.NewRequest("POST", "/", bytes.NewReader(jsonQuery))
	r = r.WithContext(context.WithValue(r.Context(), "Id", "1"))
	r = mux.SetURLVars(r, map[string]string{"post_id": "2"})
	w := httptest.NewRecorder()
	service.AddComment(w, r)
}

func TestPostHandlerPostUpvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}

	// OK
	curUser := &user.User{
		Id:       "1",
		Username: "user",
	}
	curPost := post.Post{
		Id: "2",
	}
	postRepoMock.EXPECT().PostUpvote(curPost.Id, curUser.Id).Return()
	postRepoMock.EXPECT().GetById(curPost.Id).Return(&curPost, nil)
	r := httptest.NewRequest("POST", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), "Id", curUser.Id))
	r = mux.SetURLVars(r, map[string]string{"post_id": curPost.Id})
	w := httptest.NewRecorder()
	service.PostUpvote(w, r)
}

func TestPostHandlerPostUnvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}

	// OK
	curUser := &user.User{
		Id:       "1",
		Username: "user",
	}
	curPost := post.Post{
		Id: "2",
	}
	postRepoMock.EXPECT().PostUnvote(curPost.Id, curUser.Id).Return()
	postRepoMock.EXPECT().GetById(curPost.Id).Return(&curPost, nil)
	r := httptest.NewRequest("POST", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), "Id", curUser.Id))
	r = mux.SetURLVars(r, map[string]string{"post_id": curPost.Id})
	w := httptest.NewRecorder()
	service.PostUnvote(w, r)
}

func TestPostHandlerPostDownvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}

	// OK
	curUser := &user.User{
		Id:       "1",
		Username: "user",
	}
	curPost := post.Post{
		Id: "2",
	}
	postRepoMock.EXPECT().PostDownvote(curPost.Id, curUser.Id).Return()
	postRepoMock.EXPECT().GetById(curPost.Id).Return(&curPost, nil)
	r := httptest.NewRequest("POST", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), "Id", curUser.Id))
	r = mux.SetURLVars(r, map[string]string{"post_id": curPost.Id})
	w := httptest.NewRecorder()
	service.PostDownvote(w, r)
}

func TestPostHandlerCommentDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}

	// OK
	curUser := &user.User{
		Id:       "1",
		Username: "user",
	}
	curPost := post.Post{
		Id: "2",
	}
	curComment := comment.Comment{
		Id: "3",
	}
	postRepoMock.EXPECT().CommentDelete(curPost.Id, curComment.Id, curUser.Id).Return()
	postRepoMock.EXPECT().GetById(curPost.Id).Return(&curPost, nil)
	r := httptest.NewRequest("POST", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), "Id", curUser.Id))
	r = mux.SetURLVars(r, map[string]string{"post_id": curPost.Id, "comment_id": curComment.Id})
	w := httptest.NewRecorder()
	service.CommentDelete(w, r)
}

func TestPostHandlerPostDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepoMock := NewMockPostRepo(ctrl)
	userRepoMock := NewMockUserRepo(ctrl)
	service := &PostHandler{
		UserRepo: userRepoMock,
		PostRepo: postRepoMock,
	}

	// User didn't publish the post
	curUser := &user.User{
		Id:       "1",
		Username: "user",
	}
	curPost := post.Post{
		Id: "2",
		Author: &user.User{
			Id: "3",
		},
	}
	postRepoMock.EXPECT().GetById(curPost.Id).Return(&curPost, nil)
	r := httptest.NewRequest("POST", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), "Id", curUser.Id))
	r = mux.SetURLVars(r, map[string]string{"post_id": curPost.Id})
	w := httptest.NewRecorder()
	service.PostDelete(w, r)

	// OK
	curUser = &user.User{
		Id:       "1",
		Username: "user",
	}
	curPost = post.Post{
		Id: "2",
		Author: &user.User{
			Id: "1",
		},
	}
	postRepoMock.EXPECT().GetById(curPost.Id).Return(&curPost, nil)
	postRepoMock.EXPECT().PostDelete(curPost.Id).Return()
	r = httptest.NewRequest("POST", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), "Id", curUser.Id))
	r = mux.SetURLVars(r, map[string]string{"post_id": curPost.Id})
	w = httptest.NewRecorder()
	service.PostDelete(w, r)
}

// go test -v -coverprofile=handler.out && go tool cover -html=handler.out -o handler.html && rm handler.out
