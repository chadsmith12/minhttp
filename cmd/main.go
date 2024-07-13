package main

import (
	"fmt"
	"net"
	"os"

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
    
    req, err := minhttp.Read(conn)
    if err != nil {
        fmt.Printf("Error reading http request: %s\n", err.Error())
        os.Exit(1)
    }

    if req.Path != "/" {
        minhttp.WriteNotFound(conn)
        //conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
    } else {
        minhttp.WriteOk(conn)
        //conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
    }
}
