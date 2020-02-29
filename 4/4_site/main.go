package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// GET - получение

// POST - добавление новых данных
// PUT - изменение данных
// DELETE - удаление

// HEAD
// PATCH
// OPTIONS

type User struct {
	ID       int
	Login    string
	Password string
}

func main() {

	users := map[string]*User{
		"test": &User{
			ID:       1,
			Login:    "test",
			Password: "test",
		},
	}

	sessions := map[string]*User{
		"tokenknsjkdfklsdf": users["test"],
	}

	mu := &sync.Mutex{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("session_id")
		if err != nil {
			w.Write([]byte("error" + err.Error()))
			return
		}
		mu.Lock()
		user, ok := sessions[sessionID.Value]
		mu.Unlock()
		if ok {
			w.Write([]byte(`
			<!doctype html>
			<html>
			<body>
			` + user.Login + `
				<form action="/logout" method="post">
					<button type="submit">Logout</button>
				</form>

			</body>
			</html>
			`))
			return
		}

		w.Write([]byte("no session"))
	})

	// GET
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

	})

	// POST ?login=dmitry&password=1234
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

	})

	// DELETE
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Write([]byte("wrong method"))
		}

		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Printf("cookie err: %s", err)
		}

		cookie.Expires = time.Now().Add(-1)

		http.SetCookie(w, cookie)

		mu.Lock()
		delete(sessions, cookie.Value)
		mu.Unlock()

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}
