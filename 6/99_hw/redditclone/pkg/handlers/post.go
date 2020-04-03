package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/comment"
	"redditclone/pkg/post"
	"redditclone/pkg/user"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostHandler struct {
	PostRepo *post.PostRepo
	UserRepo *user.UserRepo
	Logger   *zap.SugaredLogger
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
	var newPost post.Post
	read, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(read, &newPost)
	newPost.Author = curUser
	newPost.Votes = []post.Vote{
		post.Vote{
			UserId: newPost.Author.Id,
			Vote:   1,
		},
	}
	rnd := make([]byte, 16)
	rand.Read(rnd)
	newPost.Id = base64.URLEncoding.EncodeToString(rnd)
	newPost.Score = 1
	newPost.Views = 0
	newPost.Created = time.Now()
	newPost.UpvotePercentage = 100
	newPost.Comments = make([]comment.Comment, 0, 10)
	u.PostRepo.AddPost(&newPost)
	res, _ := json.Marshal(newPost)
	w.Write(res)
}

func (u *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	var unm map[string]string
	read, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(read, &unm)
	curUser, _ := u.UserRepo.GetUserById(r.Context().Value("Id").(string))
	rnd := make([]byte, 16)
	rand.Read(rnd)
	newComment := comment.Comment{
		Created: time.Now(),
		Author:  curUser,
		Body:    unm["comment"],
		Id:      base64.URLEncoding.EncodeToString(rnd),
	}
	u.PostRepo.AddCommentToPost(postId, newComment)
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
