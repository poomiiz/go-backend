package main

import (
  "fmt"
  "net/http"
  "os"
)

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from go-backend!")
  })
  fmt.Println("Server listening on port", port)
  if err := http.ListenAndServe(":"+port, nil); err != nil {
    panic(err)
  }
}
