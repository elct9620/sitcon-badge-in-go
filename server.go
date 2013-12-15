package main

import(
  "github.com/codegangsta/martini"
  "net/http"
)

func main() {
  m := martini.Classic()
  m.Get("/", func () string {
    return "Hello World"
  })
  m.Get("/badge", func (res http.ResponseWriter, req *http.Request) string {
    return "Host: " + req.URL.Host
  })
  m.Use(martini.Static("assets"))
  m.Run()
}
