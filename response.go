package minhttp

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"slices"
)

type HttpResponseBuilder struct {
    statusCode int
    statusMessage string
    conn io.Writer
    headers *HeadersCollection
    request *HttpRequest
    body []byte
}

func NewResponse(writer io.Writer, request *HttpRequest) *HttpResponseBuilder {
    return &HttpResponseBuilder{
    	statusCode:    200,
    	statusMessage: "OK",
    	conn:          writer,
    	headers:       NewHeadersCollection(),
	request:       request,
    	body:          []byte{},
    } 
}

func (responseBuilder *HttpResponseBuilder) WithStatus(statusCode int, statusMessage string) *HttpResponseBuilder {
    responseBuilder.statusCode = statusCode
    responseBuilder.statusMessage = statusMessage

    return responseBuilder
}

func (responseBuilder *HttpResponseBuilder) WithContent(contentType string, content []byte) *HttpResponseBuilder {
    // if we support gzip, use it
    if slices.Contains(responseBuilder.request.AcceptEncoding, "gzip") {
	return responseBuilder.WithGzip(contentType, content)
    }

    responseBuilder.headers.Add("Content-Type", contentType)
    responseBuilder.headers.Add("Content-Lenth", fmt.Sprint(len(content)))
    responseBuilder.body = content
    return responseBuilder
}

func (responseBuilder *HttpResponseBuilder) WithGzip(contentType string, content []byte) *HttpResponseBuilder {
    var gzipBuffer bytes.Buffer
    gzipWriter := gzip.NewWriter(&gzipBuffer)
    gzipWriter.Write(content)
    gzipWriter.Close()

    responseBuilder.headers.Add("Content-Encoding", "gzip")
    responseBuilder.headers.Add("Content-Type", contentType)
    responseBuilder.headers.Add("Content-Lenth", fmt.Sprint(len(gzipBuffer.Bytes())))
    responseBuilder.body = gzipBuffer.Bytes()
    return responseBuilder
}

func (responseBuilder *HttpResponseBuilder) WithHeader(key, value string) *HttpResponseBuilder {
    responseBuilder.headers.Add(key, value)
    return responseBuilder
}

func (responseBuilder *HttpResponseBuilder) Write() {
    writeStatus(responseBuilder.conn, responseBuilder.statusCode, responseBuilder.statusMessage)
    for header, value := range responseBuilder.headers.rawHeaders {
	fmt.Fprintf(responseBuilder.conn, "%s: %s\r\n", header, value)
    }
    fmt.Fprint(responseBuilder.conn, "\r\n")

    if len(responseBuilder.body) > 0 {
	responseBuilder.conn.Write(responseBuilder.body)
    }
}

func WriteText(builder *HttpResponseBuilder, text string) {
    builder.WithContent("text/plain", []byte(text))
    builder.Write()
}

func WriteOctetStream(builder *HttpResponseBuilder, content []byte) {
    builder.WithContent("application/octet-stream", content)
    builder.Write()
}

func WriteOk(builder *HttpResponseBuilder) {
    builder.Write()
}

func WriteCreated(builder *HttpResponseBuilder) {
    builder.WithStatus(201, "Created")
    builder.Write()
}

func WriteNotFound(builder *HttpResponseBuilder) {
    builder.WithStatus(404, "Not Found")
    builder.Write()
}

func WriteBadRequest(builder *HttpResponseBuilder) {
    builder.WithStatus(400, "Bad Request")
    builder.Write()
}

func WriteInternalServerError(builder *HttpResponseBuilder) {
    builder.WithStatus(500, "Internal Server Error")
    builder.Write()
}

func writeStatus(writer io.Writer, statusCode int, statusText string) {
    statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
    writer.Write([]byte(statusLine))
}
