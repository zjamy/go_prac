package main

import (
  "net/http"
  "os/signal"
  "fmt"
  "io"
  "errgroup"
)

func StartHttpServer(srv *http.Server) error {
  http.HandleFunc("/hello", HelloServer2)
  fmt.Println("http server start")
  err := srv.ListenAndServe()
  return err
}

func HelloServer2(iw http.ResponseWriter, req *http.Request) {
  io.WriteString(w, "Hello, World!\n")
}

func main() {
  ctx := context.Background()
  ctx, cancel := context.WithCancel(ctx)
  group, errCtx := errgroup.WithContext(ctx)
  srv := &http.Server{Addr: "8000"}

  group.GO(func() error {
    return StartHttpServer(srv)
})

  group.Go(func() error {
    <-errCtx.Done()
    fmt.Println("http server stop")
    return srv.Shutdown(errCtx)
  })

  chanel := make(chan os.Signal, 1)
  signal.Notify(chanel)

  group.Go(func() error {
    for {
      select {
      case <-errCtx.Done():
        return errCtx.Err()
      case <-chanel:
        cancel()

      }
    }
    return nil
  })

  if err := group.Wait(); err != nil {
    fmt.Println("group error: ", err)
  }
  fmt.Println("all group done!")
}
