package main
import (
  "net/http"
  "log"
)


type myHandler struct{
  content string
}

func (handler *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(handler.content))
}
func SayBye(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("bye bye, this is v2 httpServer"))
}

func main() {
  http.Handle("/", &myHandler{content: "this is v2 httpServer"})
  http.HandleFunc("/bye",SayBye)

  log.Println("Starting v2 httpserver")
  err := http.ListenAndServe(":8000",nil)
  if err != nil {
    panic(err)
    //panic("Listen and Serve Err!")
  }
}
