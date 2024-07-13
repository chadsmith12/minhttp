package minhttp

import (
	"fmt"
	"io"
)

func WriteText(writer io.Writer, text string) {
    writeStatus(writer, 200, "OK")
    writeContentHeaders(writer, "text/plain", []byte(text))
    writer.Write([]byte("\r\n"))
    writer.Write([]byte(text))
}

func WriteOk(writer io.Writer) {
    writeStatus(writer, 200, "OK")
    writer.Write([]byte("\r\n"))
}

func WriteNotFound(writer io.Writer) {
    writeStatus(writer, 404, "Not Found")
    writer.Write([]byte("\r\n"))
}

func writeContentHeaders(writer io.Writer, contentType string, content []byte) {
    contentHeader := fmt.Sprintf("Content-Type: %s\r\n", contentType)
    contentLength := fmt.Sprintf("Content-Length: %d\r\n", len(content))
    writer.Write([]byte(contentHeader))
    writer.Write([]byte(contentLength))
}

func writeStatus(writer io.Writer, statusCode int, statusText string) {
    statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
    writer.Write([]byte(statusLine))
}

func writeNotFound(writer io.Writer) {
    writer.Write([]byte("HTTP/1.1 404 Not Found"))
}

