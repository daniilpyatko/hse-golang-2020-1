package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/comment"
	"redditclone/pkg/post"
	"redditclone/pkg/user"

	"github.com/gorilla/mux"
)

// mockgen -source=post.go -destination=post_mock.go -package=handlers PostRepo

type PostRepo interface {
	GetAll() ([]*post.Post, error)
	GetByCategory(category string) ([]*post.Post, error)
	GetById(id string) (*post.Post, error)
	GetByUsername(username string) ([]*post.Post, error)
	AddPost(post *post.Post)
	AddNewPost(user *user.User, newPost *post.Post)
	AddCommentToPost(postId string, curUser *user.User, curComment comment.Comment)
	PostUpvote(postId, curUserId string)
	PostDownvote(postId, curUserId string)
	PostUnvote(postId, curUserId string)
	CommentDelete(postId, commentId, curUserId string)
	PostDelete(postId string)
}

type PostHandler struct {
	PostRepo PostRepo
	UserRepo UserRepo
}

func (u *PostHandler) AllPosts(w http.ResponseWriter, r *http.Request) {
	cur, _ := u.PostRepo.GetAll()
	res, _ := json.Marshal(cur)
	w.Write(res)
}

func (u *PostHandler) PostsByCategory(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)["category"]
	cur, _ := u.PostRepo.GetByCategory(category)
	res, _ := json.Marshal(cur)
	w.Write(res)
}

func (u *PostHandler) PostsById(w http.ResponseWriter, r *http.Request) {
	Id := mux.Vars(r)["post_id"]
	cur, _ := u.PostRepo.GetById(Id)
	cur.Views += 1
	u.PostRepo.PostDelete(Id)
	u.PostRepo.AddPost(cur)
	res, _ := json.Marshal(cur)
	w.Write(res)
}

func (u *PostHandler) PostsByUsername(w http.ResponseWriter, r *http.Request) {
	Username := mux.Vars(r)["username"]
	cur, _ := u.PostRepo.GetByUsername(Username)
	res, _ := json.Marshal(cur)
	w.Write(res)
}

func (u *PostHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	curUser, err := u.UserRepo.GetUserById(r.Context().Value("Id").(string))
	if err != nil {
		return
	}
	newPost := post.Post{}
	read, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(read, &newPost)
	u.PostRepo.AddNewPost(curUser, &newPost)
	res, _ := json.Marshal(newPost)
	w.Write(res)
}

func (u *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	var unm map[string]string
	read, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(read, &unm)
	curUser, _ := u.UserRepo.GetUserById(r.Context().Value("Id").(string))
	newComment := comment.Comment{
		Body: unm["comment"],
	}
	u.PostRepo.AddCommentToPost(postId, curUser, newComment)
	resPost, _ := u.PostRepo.GetById(postId)
	res, _ := json.Marshal(resPost)
	w.Write(res)
}

func (u *PostHandler) PostUpvote(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	curUserId, _ := r.Context().Value("Id").(string)
	u.PostRepo.PostUpvote(postId, curUserId)
	Post, _ := u.PostRepo.GetById(postId)
	res, _ := json.Marshal(Post)
	w.Write(res)
}

func (u *PostHandler) PostDownvote(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	curUserId, _ := r.Context().Value("Id").(string)
	u.PostRepo.PostDownvote(postId, curUserId)
	Post, _ := u.PostRepo.GetById(postId)
	res, _ := json.Marshal(Post)
	w.Write(res)
}

func (u *PostHandler) PostUnvote(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	curUserId, _ := r.Context().Value("Id").(string)
	u.PostRepo.PostUnvote(postId, curUserId)
	Post, _ := u.PostRepo.GetById(postId)
	res, _ := json.Marshal(Post)
	w.Write(res)
}

func (u *PostHandler) CommentDelete(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	commentId := mux.Vars(r)["comment_id"]
	curUserId, _ := r.Context().Value("Id").(string)
	u.PostRepo.CommentDelete(postId, commentId, curUserId)
	Post, _ := u.PostRepo.GetById(postId)
	res, _ := json.Marshal(Post)
	w.Write(res)
}

func (u *PostHandler) PostDelete(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	curUserId, _ := r.Context().Value("Id").(string)
	Post, _ := u.PostRepo.GetById(postId)
	if Post.Author.Id == curUserId {
		u.PostRepo.PostDelete(postId)
		message := map[string]string{
			"message": "success",
		}
		jsonmessage, _ := json.Marshal(message)
		w.Write(jsonmessage)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
