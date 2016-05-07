package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "345 Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	fmt.Println("!!Hello from the outside!!~~")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
