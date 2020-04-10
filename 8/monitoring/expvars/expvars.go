package main

import (
	"expvar"
	"fmt"
	"net/http"
)

var (
	hits = expvar.NewMap("hits")
	i    = 0
)

func handler(w http.ResponseWriter, r *http.Request) {
	hits.Add(r.URL.String(), 1)

	fmt.Println("hit to" + r.URL.String())

	w.Write([]byte("expvar increased"))
	i++
}

func main() {
	http.HandleFunc("/", handler)

	expvar.Publish("mystat", expvar.Func(func() interface{} {
		return map[string]int{
			"test":  100500,
			"value": i,
		}
	}))

	fmt.Println("starting server at :8082")
	http.ListenAndServe(":8082", nil)
}
