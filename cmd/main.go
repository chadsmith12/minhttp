package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/chadsmith12/minhttp"
)

func main() {
    listener, err := net.Listen("tcp", "0.0.0.0:4221")
    if err != nil {
        fmt.Println("Failed to bind to port 4221")
        os.Exit(1)
    }

    conn, err := listener.Accept()
    if err != nil {
        fmt.Printf("Error accepting connection: %v\n", err.Error())
        os.Exit(1)
    }
    
    req, err := minhttp.ReadRequest(conn)
    if err != nil {
        fmt.Printf("Error reading http request: %s\n", err.Error())
        os.Exit(1)
    }

    if req.Path == "/" {
        minhttp.WriteOk(conn)
    } else if strings.HasPrefix(req.Path, "/echo") {
        components := strings.Split(req.Path, "/echo/")
        minhttp.WriteText(conn, components[1]) 
    } else if strings.HasPrefix(req.Path, "/user-agent") {
        userAgent := req.Headers.UserAgent()
        minhttp.WriteText(conn, userAgent)
    } else {
        minhttp.WriteNotFound(conn)
    }
}
