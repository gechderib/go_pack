package main

import (
	"fmt"
	"net/http"
	"time"
)

func handleHello(w http.ResponseWriter, r *http.Request) {

	fmt.Println("The request method ", r.Method)
	fmt.Println("The request header is ", r.Header)
	fmt.Println("The request url is ", r.URL)
	fmt.Println("The request body is ", r.Body)

	name := r.URL.Query().Get("name")
	fmt.Println(name)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("Hello Go Backend!"))

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The request method ", r.Method)
	fmt.Println("The request header is ", r.Header)
	fmt.Println("The request url is ", r.URL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("Welcome to the Home Page!"))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request: 5555555555 ", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})

}

func main() {

	http.Handle("/", loggingMiddleware(http.HandlerFunc(handleHome)))
	http.Handle("/hello", loggingMiddleware(http.HandlerFunc(handleHello)))

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("The request method ", r.Method)
		fmt.Println("The request header is ", r.Header)
		fmt.Println("The request url is ", r.URL)

		w.Write([]byte("The status is Confirmed"))
	})

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("The request method ", r.Method)
		fmt.Println("The request header is ", r.Header)
		fmt.Println("The request url is ", r.URL)

		w.Write([]byte(fmt.Sprintf("The current time is %s", time.Now().Format(time.RFC1123))))
	})

	http.ListenAndServe(":8080", nil)
}
