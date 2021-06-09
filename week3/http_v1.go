package main
import (
  "fmt"
  "net/http"
)

func SayBye(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Bye Bye!")
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello World!")
}

func main () {
  http.HandleFunc("/", HelloHandler)
  http.HandleFunc("/bye", SayBye)
  http.ListenAndServe(":8000",nil)
}
