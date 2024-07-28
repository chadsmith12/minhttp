package main

import (
	"strings"

	"github.com/chadsmith12/minhttp"
)

func main() {
    input := "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n"
    reader := strings.NewReader(input)
    minhttp.ReadRequest(reader)
}
