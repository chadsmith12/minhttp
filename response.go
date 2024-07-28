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

func WriteOctetStream(writer io.Writer, content []byte) {
    writeStatus(writer, 200, "OK")
    writeContentHeaders(writer, "application/octet-stream", content)
    writer.Write([]byte("\r\n"))
    writer.Write(content)
}

func WriteOk(writer io.Writer) {
    writeStatus(writer, 200, "OK")
    writer.Write([]byte("\r\n"))
}

func WriteCreated(writer io.Writer) {
    writeStatus(writer, 201, "Created")
    writer.Write([]byte("\r\n"))
}

func WriteNotFound(writer io.Writer) {
    writeStatus(writer, 404, "Not Found")
    writer.Write([]byte("\r\n"))
}

func WriteBadRequest(writer io.Writer) {
    writeStatus(writer, 400, "Bad Request")
    writer.Write([]byte("\r\n"))
}

func WriteInernalServerError(writer io.Writer, err string) {
   writeStatus(writer, 500, fmt.Sprintf("Internal Server Error - %s", err)) 
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

