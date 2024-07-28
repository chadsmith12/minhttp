package main

import (
	"fmt"
	"io"

	"github.com/chadsmith12/minhttp"
)

func main() {
    router := minhttp.NewRouter()
    hander := func(hr minhttp.HttpRequest, w io.Writer) error { return nil}
    router.MapGet("/users", hander)
    router.MapGet("/users/list", hander)
    router.MapGet("/users/{id}", hander)
    router.MapGet("/users/{id}/profile", hander)
    router.MapGet("/echo/hello-world", hander)
    router.MapGet("/user-agent", hander)
    router.MapGet("/temp/{str}", hander)
    router.MapGet("/users/{id}/messages", hander)
    
    fmt.Println(router.MatchRoute("GET", "/user-agent"))
    fmt.Println(router.MatchRoute("GET", "/users"))
    fmt.Println(router.MatchRoute("GET","/users/10"))
    fmt.Println(router.MatchRoute("GET","/users/10/profile"))
    fmt.Println(router.MatchRoute("GET","/temp/hello"))
    fmt.Println(router.MatchRoute("GET","/profile"))
    fmt.Println(router.MatchRoute("GET","/users/list"))
    fmt.Println(router.MatchRoute("GET","/users/10/messages"))
}
