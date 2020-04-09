package main

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/post"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	dsn := "root:My111111@tcp(localhost:3306)/db?"
	dsn += "&charset=utf8"
	db, err := sql.Open("mysql", dsn)
	db.Exec("CREATE TABLE IF NOT EXISTS sessions (`id` VARCHAR(255), `userid` VARCHAR(255));")
	db.Exec("CREATE TABLE IF NOT EXISTS users (`id` VARCHAR(255), `username` VARCHAR(255), `password` VARCHAR(255));")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	connectionURI := "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	err = client.Connect(context.TODO())
	collection := client.Database("db").Collection("posts")

	tmpl := template.Must(template.ParseGlob("./template/*html"))
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	postRepo := post.NewRepo(collection)
	userRepo := user.NewRepo(db)
	sessionManager := session.NewSessionManager(db)
	postHandler := handlers.PostHandler{
		PostRepo: postRepo,
		UserRepo: userRepo,
	}
	userHandler := handlers.UserHandler{
		UserRepo:       userRepo,
		SessionManager: sessionManager,
	}
	rMain := mux.NewRouter()
	rWithAuth := mux.NewRouter()
	rMain.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	rMain.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./template/static/"))))

	rMain.HandleFunc("/api/posts/", postHandler.AllPosts).Methods("GET")
	rMain.HandleFunc("/api/posts/{category}", postHandler.PostsByCategory).Methods("GET")
	rMain.HandleFunc("/api/post/{post_id}", postHandler.PostsById).Methods("GET")
	rMain.HandleFunc("/api/user/{username}", postHandler.PostsByUsername).Methods("GET")
	rMain.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	rMain.HandleFunc("/api/register", userHandler.Register).Methods("POST")

	rWithAuth.HandleFunc("/api/posts", postHandler.NewPost).Methods("POST")
	rWithAuth.HandleFunc("/api/post/{post_id}", postHandler.AddComment).Methods("POST")
	rWithAuth.HandleFunc("/api/post/{post_id}/upvote", postHandler.PostUpvote).Methods("GET")
	rWithAuth.HandleFunc("/api/post/{post_id}/downvote", postHandler.PostDownvote).Methods("GET")
	rWithAuth.HandleFunc("/api/post/{post_id}/unvote", postHandler.PostUnvote).Methods("GET")
	rWithAuth.HandleFunc("/api/post/{post_id}/{comment_id}", postHandler.CommentDelete).Methods("DELETE")
	rWithAuth.HandleFunc("/api/post/{post_id}", postHandler.PostDelete).Methods("DELETE")

	rWithAuthMiddleware := middleware.Auth(sessionManager, logger, rWithAuth)

	rMain.Handle("/api/posts", rWithAuthMiddleware).Methods("POST")
	rMain.Handle("/api/post/{post_id}", rWithAuthMiddleware).Methods("POST")
	rMain.Handle("/api/post/{post_id}/upvote", rWithAuthMiddleware).Methods("GET")
	rMain.Handle("/api/post/{post_id}/downvote", rWithAuthMiddleware).Methods("GET")
	rMain.Handle("/api/post/{post_id}/unvote", rWithAuthMiddleware).Methods("GET")
	rMain.Handle("/api/post/{post_id}/{comment_id}", rWithAuthMiddleware).Methods("DELETE")
	rMain.Handle("/api/post/{post_id}", rWithAuthMiddleware).Methods("DELETE")

	rMainMiddleware := middleware.AccessLog(logger, rMain)
	rMainMiddleware = middleware.Panic(logger, rMainMiddleware)
	addr := ":8080"
	http.ListenAndServe(addr, rMainMiddleware)
}
