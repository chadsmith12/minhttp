package main

import (
	"fmt"
	"github.com/chadsmith12/minhttp"
)

func main() {
    router := minhttp.NewRouter()
    router.MapGet("/users")
    router.MapGet("/users/list")
    router.MapGet("/users/{id}")
    router.MapGet("/users/{id}/profile")
    router.MapGet("/echo/hello-world")
    router.MapGet("/user-agent")
    router.MapGet("/temp/{str}")
    router.MapGet("/users/{id}/messages")
    
    fmt.Println(router.MatchRoute("/user-agent"))
    fmt.Println(router.MatchRoute("/users"))
    fmt.Println(router.MatchRoute("/users/10"))
    fmt.Println(router.MatchRoute("/users/10/profile"))
    fmt.Println(router.MatchRoute("/temp/hello"))
    fmt.Println(router.MatchRoute("/profile"))
    fmt.Println(router.MatchRoute("/users/list"))
    fmt.Println(router.MatchRoute("/users/10/messages"))
}
