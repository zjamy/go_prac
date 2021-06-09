package main

import (
    "log"
    "os"
    "os/signal"
    "time"
    "net/http"
)

var server *http.Server
func main() {
  quit := make(chan os.Signal)
  signal.Notify(quit, os.Interrupt)

  mux := http.NewServeMux()
  mux.Handle("/", &myHandler{})
  mux.HandleFunc("/bye", sayBye)

  server = &http.Server{
    Addr:         ":8000",
    WriteTimeout: time.Second * 4,
    Handler:      mux,
  }

  go func() {
    <-quit
    if err := server.Close(); err != nil {
        log.Fatal("Close server:", err)
    }
  }()

  log.Println("Starting v3 httpserver")
  err := server.ListenAndServe()
  if err != nil {
    if err == http.ErrServerClosed {
      log.Fatal("Server closed under request")
    } else {
      log.Fatal("Server closed unexpected", err)
    }
  }
  log.Fatal("Server exited")

}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("this is version 3"))
}

func sayBye(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("bye bye ,shutdown the server"))
  err := server.Shutdown(nil)
  if err != nil {
    log.Fatal([]byte("shutdown the server err"))
  }
}
