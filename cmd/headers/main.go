package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chadsmith12/minhttp"
)

func main() {
    input2 := "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n"
    reader2 := strings.NewReader(input2)
    _, err := minhttp.ReadRequest(reader2)
    if err != nil {
        fmt.Println("failed to parsed request")
        os.Exit(1)
    }
}
