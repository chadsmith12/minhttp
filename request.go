package minhttp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
)

type HttpRequest struct {
    Method string
    Path string
    Version string
    Headers *HeadersCollection
    Params map[string]string
    ContentLength int64
    AcceptEncoding []string
    Body io.Reader
}

type HeadersCollection struct {
    rawHeaders map[string]string
}

func NewHeadersCollection() *HeadersCollection {
   return &HeadersCollection{
   	rawHeaders: map[string]string{},
   } 
}

func (headers *HeadersCollection) AddRaw(headerBytes []byte) {
    key, value, found := bytes.Cut(headerBytes, []byte{':'})
    if !found {
        return
    }
    keyText := strings.TrimSpace(string(key))
    valueText := strings.TrimSpace(string(value))
    
    current, ok := headers.Get(keyText)
    if ok {
        current = fmt.Sprintf("%s,%s", current, valueText)
        return
    }

    headers.rawHeaders[keyText] = valueText
}

func (headers *HeadersCollection) Add(header string, value string) {
    headers.rawHeaders[header] = value
}

func (headers *HeadersCollection) Get(header string) (string, bool) {
    value, ok := headers.rawHeaders[header]

    return value, ok
}

func (headers *HeadersCollection) Remove(header string) {
   delete(headers.rawHeaders, header) 
}

func (headers *HeadersCollection) Len() int {
    return len(headers.rawHeaders)
}

func (headers *HeadersCollection) UserAgent() string {
    userAgent, _ := headers.Get("User-Agent")

    return userAgent
}

func (headers *HeadersCollection) ContentLength() int64 {
    contentLengthStr, ok := headers.Get("Content-Length")
    if !ok {
        return 0
    }

    contentLength, err := strconv.Atoi(contentLengthStr)
    if err != nil {
        return 0
    }

    return int64(contentLength)
}

func (headers *HeadersCollection) AcceptEncoding() []string {
    encodings, ok := headers.Get("Accept-Encoding")
    if !ok {
        return []string{}
    }

    if len(encodings) == 0 {
        return []string{}
    }

    acceptedEncodings := strings.Split(encodings, ", ")
    return acceptedEncodings 
}

func ReadRequest(reader io.Reader) (HttpRequest, error) {
    bufReader := bufio.NewReader(reader)
    textReader := textproto.NewReader(bufReader)
    statusLine, err := textReader.ReadLine()
    method, target, version, err := decodeStatusLine(statusLine)
    if err != nil {
        return HttpRequest{}, err
    }

    readHeaders, err := textReader.ReadMIMEHeader()
    if err != nil {
        return HttpRequest{}, err
    }
    headers := NewHeadersCollection() 
    for key, value := range readHeaders {
        headers.Add(key, value[0])
    }

    bodyReader := io.LimitReader(bufReader, headers.ContentLength())

    return HttpRequest{
    	Method:  method,
    	Path:    target,
    	Version: version,
        Headers: headers,
        Params: make(map[string]string),
        ContentLength: headers.ContentLength(),
        AcceptEncoding: headers.AcceptEncoding(),
        Body: bodyReader,
    }, nil
}


func decodeStatusLine(statusLine string) (string, string, string, error) {
    lineData := strings.Split(statusLine, " ")
    if len(lineData) < 3 {
        return "", "", "", errors.New("failed to read the status line.")
    }

    return lineData[0], lineData[1], lineData[2], nil
}

