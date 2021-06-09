package main

import (
    "fmt"
    "net/http"
)

type WelcomeHandlerStruct struct {

}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World")
}

func (*WelcomeHandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome")
}

func main () {
    mux := http.NewServeMux()
    mux.HandleFunc("/", HelloHandler)
    mux.Handle("/welcome", &WelcomeHandlerStruct{})
    http.ListenAndServe(":8080", mux)
}
