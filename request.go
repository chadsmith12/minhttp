package minhttp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

type HttpRequest struct {
    Method string
    Path string
    Version string
}

func Read(reader io.Reader) (HttpRequest, error) {
    scanner := bufio.NewScanner(reader)
    scanner.Split(split)

    found := scanner.Scan()
    if !found {
        return HttpRequest{}, errors.New("failed to find http request to read")
    }

    if scanner.Err() != nil {
        return HttpRequest{}, errors.New("failed to read http request")
    }

    fmt.Printf("Read a total of %d bytes\n", len(scanner.Bytes()))
    return DecodeRequest(scanner.Bytes())
}

func DecodeRequest(msg []byte) (HttpRequest, error) {
    requestLine := readLine(msg)
    method, target, version := decodeRequestLine(requestLine)

    return HttpRequest{
        Method: method,
        Path: target,
        Version: version,
    }, nil
}

func readLine(data []byte) []byte {
    if index := bytes.Index(data, []byte("\r\n")); index >= 0 {
        return data[:index] 
    }

    return data
}

func decodeRequestLine(requestLine []byte) (string, string, string) {
    requestData := bytes.Split(requestLine, []byte{' '}) 
    
    return string(requestData[0]), string(requestData[1]), string(requestData[2])
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
    fmt.Printf("CALLING SPLIT WITH %s\n", string(data))
    header, _, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
    if !found {
        fmt.Println("Failed to find \\r\\n\\r\\n")
        return 0, nil, nil
    }

    totalLength := len(header) + 4

    return totalLength, header, nil
}
