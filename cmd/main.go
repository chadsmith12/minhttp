package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
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
    
    req, err := readHttp(conn)
    if err != nil {
        fmt.Printf("Error reading http request: %s\n", err.Error())
        os.Exit(1)
    }

    fmt.Printf("Read the following request: %s\n", req)
    conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

func readHttp(conn net.Conn) (string, error) {
    reader := bufio.NewReader(conn)
    var buffer bytes.Buffer
    for {
        bytesRead, err := reader.ReadBytes('\r')
        if err != nil {
            if err == io.EOF {
                break
            }

            return "", err
        }
        buffer.Write(bytesRead)
    }

    return buffer.String(), nil
}

var splitFunc bufio.SplitFunc
