package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Logging logs all requests with its path and the time it took to process
func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { log.Println(r.URL.Path, time.Since(start)) }()

		f(w, r)
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		f(w, r)
	}
}

// Chain applies middlewares to an http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main() {
	http.HandleFunc("/", Chain(Hello, func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			Method("GET", Logging(h))(w, r)
		}
	}))
	http.ListenAndServe(":8080", nil)
}
