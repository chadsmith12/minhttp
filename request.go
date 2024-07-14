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
    Headers *HeadersCollection
}

type HeadersCollection struct {
    rawHeaders map[string]string
}

func NewHeadersCollection() *HeadersCollection {
   return &HeadersCollection{
   	rawHeaders: map[string]string{},
   } 
}

func (headers *HeadersCollection) Add(headerBytes []byte) {
    key, value, found := bytes.Cut(headerBytes, []byte{':'})
    if !found {
        return
    }
    keyText := string(key)
    valueText := string(value)
    
    current, ok := headers.Get(keyText)
    if ok {
        current = fmt.Sprintf("%s,%s", current, valueText)
        return
    }

    headers.rawHeaders[keyText] = valueText
}

func (headers *HeadersCollection) Get(header string) (string, bool) {
    value, ok := headers.rawHeaders[header]

    return value, ok
}

func (headers *HeadersCollection) Len() int {
    return len(headers.rawHeaders)
}

func (headers *HeadersCollection) UserAgent() string {
    userAgent, _ := headers.Get("User-Agent")

    return userAgent
}

func ReadRequest(reader io.Reader) (HttpRequest, error) {
    scanner := bufio.NewScanner(reader)
    scanner.Split(readRequest)

    found := scanner.Scan()
    if !found {
        return HttpRequest{}, errors.New("failed to find http request to read")
    }
    fmt.Printf("Read the Following Line: %s\n", scanner.Text())

    method, target, version := decodeRequestLine(scanner.Bytes())
    headers := &HeadersCollection{rawHeaders: make(map[string]string)}
    for scanner.Scan() {
        if scanner.Text() == "" {
            break
        }
        fmt.Printf("Read the Following Line from headers: %s\n", scanner.Text())
        lineRead := scanner.Bytes()
        headers.Add(lineRead)
    }
    return HttpRequest{
    	Method:  method,
    	Path:    target,
    	Version: version,
        Headers: headers,
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

func readRequest(data []byte, atEOF bool) (advance int, token []byte, err error) {
    requestData, _, found := bytes.Cut(data, []byte{'\r', '\n'})
    if !found {
        return 0, nil, nil
    }

    totalLength := len(requestData) + 2

    return totalLength, requestData, nil
}
