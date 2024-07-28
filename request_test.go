package minhttp_test

import (
	"strings"
	"testing"

	"github.com/chadsmith12/minhttp"
)

func TestBasicDecode(t *testing.T) {
    headers := minhttp.NewHeadersCollection()
    headers.AddRaw([]byte("Host: localhost:4221"))
    headers.AddRaw([]byte("User-Agent: curl/7.64.1"))
    headers.AddRaw([]byte("Accept: */*"))
    testCases := []struct {
        name string
        input string
        expected minhttp.HttpRequest
    } {
        {
        	name:     "Should read simple GET Request",
        	input:    "GET / HTTP/1.1\r\n\r\n",
        	expected: minhttp.HttpRequest{
                    Method: "GET",
                    Path: "/",
                    Version: "HTTP/1.1",
                    Headers: minhttp.NewHeadersCollection(),
                },
        },
        {
        	name:     "Should read GET Request that has headers",
        	input:    "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n",
        	expected: minhttp.HttpRequest{
                    Method: "GET",
                    Path: "/index.html",
                    Version: "HTTP/1.1",
                    Headers: headers,
                },
        },
    }

    for _, testCase := range testCases {
        t.Run(testCase.name, func(t *testing.T) {
            reader := strings.NewReader(testCase.input)
            req, err := minhttp.ReadRequest(reader)
            if err != nil {
                t.Fatalf("Unexpected error: %v", err)
            }
            if req.Method != testCase.expected.Method {
                t.Errorf("Expected method %s, got %s", testCase.expected.Method, req.Method)
            }
            if req.Path != testCase.expected.Path {
                t.Errorf("Expected path %s, got %s", testCase.expected.Path, req.Path)
            }
            if req.Version != testCase.expected.Version {
                t.Errorf("Expected version %s, got %s", testCase.expected.Version, req.Version)
            }
            if (req.Headers.Len() != testCase.expected.Headers.Len()) {
                t.Errorf("Expected to get %d headers, got %d", testCase.expected.Headers.Len(), req.Headers.Len())
            }
        })
    }
}
