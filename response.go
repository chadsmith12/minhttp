package minhttp

import (
	"fmt"
	"io"
)

type ResponseWriter interface {
    WriteStatus(statusCode int, statusMessage string)
    Write([]byte) (int, error)
    Headers() *HeadersCollection
}

type httpResponse struct {
    conn io.Writer
    request *HttpRequest
    headers *HeadersCollection
    statusCode int
    statusMessage string
    wroteHeader bool
}

func (resp *httpResponse) Headers() *HeadersCollection {
    return resp.headers
}

func (resp *httpResponse) WriteStatus(statusCode int, statusMessage string) {
    if resp.wroteHeader {
	return
    }

    resp.statusCode = statusCode
    resp.statusMessage = statusMessage
    resp.wroteHeader = true

    writeStatus(resp.conn, resp.statusCode, resp.statusMessage)
}

func (resp *httpResponse) Write(data []byte) (int, error) {
    resp.finishRequest()
    fmt.Fprint(resp.conn, "\r\n")
    return resp.conn.Write(data)
}

func (resp *httpResponse) finishRequest() {
    if !resp.wroteHeader {
	resp.WriteStatus(200, "OK")
    }
    
    for header, value := range resp.headers.rawHeaders {
	fmt.Fprintf(resp.conn, "%s: %s\r\n", header, value)
    }
}

func WriteText(writer ResponseWriter, text string) {
    writer.Headers().Add("Content-Type", "text/plain")
    writer.Headers().Add("Content-Length", fmt.Sprint(len(text)))
    writer.Write([]byte(text))
}

func WriteOctetStream(writer ResponseWriter, content []byte) {
    writer.Headers().Add("Content-Type", "text/plain")
    writer.Headers().Add("Content-Length", fmt.Sprint(len(content)))
    writer.Write(content)
}

func WriteOk(writer ResponseWriter) {
    writer.Write([]byte{})
}

func WriteCreated(writer ResponseWriter) {
    writer.WriteStatus(201, "Created")
    writer.Write([]byte{})
}

func WriteNotFound(writer ResponseWriter) {
    writer.WriteStatus(404, "Not Found")
    writer.Write([]byte{})
}

func WriteBadRequest(writer ResponseWriter) {
    writer.WriteStatus(400, "Bad Request")
    writer.Write([]byte{})
}

func WriteInernalServerError(writer ResponseWriter, err string) {
    writer.WriteStatus(500, fmt.Sprintf("Internal Server Error - %s", err)) 
    writer.Write([]byte{})
}

func writeStatus(writer io.Writer, statusCode int, statusText string) {
    statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
    writer.Write([]byte(statusLine))
}
