# minhttp

*NOTE: THIS SHOULD NOT BE USED IN PRODUCTION.*

Minhttp is a minimal HTTP Server written in Go from scratch. This doesn't use Go's http standard library server nor use any other library.
The main goal of this was a learning experience and learning the basics of the HTTP 1.1 protocol. This does not fully implement the HTTP 1.1 RFC
and only covers a subset of it. 

As more HTTP features were added code was refactored to try to organize it in a way to try to create a standard http application framework and also look similar to Go's `net/http` package.

## Basics

Using `minhttp` is similar to how you would use Go's `net/http` package but adds a some nice helpers that Go's standard library doesn't provide in their http package.

Starting a basic http server looks like the following:

```go
package main

import (
	"fmt"
	"os"

	"github.com/chadsmith12/minhttp"
)

func main() {
    minhttp.MapGet("/", func(hr minhttp.HttpRequest, w *minhttp.HttpResponseBuilder) error {
        minhttp.WriteOk(w)
        return nil
    })

    minhttp.MapGet("/echo/{str}", func(hr minhttp.HttpRequest, w *minhttp.HttpResponseBuilder) error {
        echo := hr.Params["str"]
        minhttp.WriteText(w, echo)	
        return nil
    })

    err := minhttp.ListenAndServe("0.0.0.0:4221")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```
This example creates two HTTP Endpoints of "/" and "/echo{str}" and then starts up a server on port 4221. It also uses the default router when the serer is started (similar to `net/http`)

Unlike the standard library where you now make the path template with the HTTP method embedded in it, `minhttp` separates out the ways you map your HTTP Methods.
The next convient thing is the `ListenAndServe` method. The `ListenAndServe` method automatically listens to the signal from the terminal to terminate the application.
So the application will always attempt to terminate cleanly.

## Implementation

The following gives the basic implementation of the different parts of the HTTP Server and how it was tackled in Go.

### Server

Everything starts with TCP. A basic wrapper about a TCP Connection is written to easily spin up the TCP Server for the HTTP Server.
For every TCP Connection it is handled on a separate go routine. This is added to a wait group so that way if the server is stopped during a request, it will attempt to wait for the request to be finished before it closes the server out fully.

The HTTP Server will automatically create a server with a default router when you call `ListenAndServe`. This gives you a quick way to get up and running without having to define a lot of different things for the serve. This gives you a quick way to get up and running without having to define a lot of different things for the server.

The TCP Server will accept the connection and start the go routine, and then call a handler that the HTTP Server provides it. This handler is what actually attempts to process the request as a HTTP Request.

### HTTP Requests

The server will parse the HTTP Request using Go's standard `textproto` library to allow for a quick way to read in the request data. If during the parsing of this request
an error occurs, then the handler will automatically write back and HTTP Response error to the connection.

The request will parse out the following details:

* Method - The HTTP Method Used. 
* Path - The path requested
* Version - HTTP version used (should always be 1.1 because this is all we support)
* Headers - HTTP Headers read
* Params - Params in the request
* ContentLength - Quick access to the content-length, if provided
* AcceptEncoding - quick access to the encoding's accepted, if provided
* Body - Access to read the body of the http request

After the request is read, found to be correct, and created, then it will attempt to match the path requested to the routes defined on the server.
If a route was found, that routes handler will be called.
If a route was not found, the server will write back a 404 Not Found response to the connection.

### HTTP Routes

A basic routing engine is implemented that supports two basic types of routes:

* Static Route - A route like `/hello-world`
* Wild Card Route - A route like `/echo/{str}`

The routing engine is implemented using the Tries data structure. When you map a route template, the routing engine will parse each segment of the route and turn each segment
into a `RouteNode`. A `RouteNode` is just a node on the Trie data structure and stores when a full route has been created.
The router itself stores the root node for the trie for it start it's traversal.

The Router has the ability to take a request from the server and attempt to match that request to some route on the router's trie.
It will preform a depth-first-search attempting to match each segment of the path provided in the request. When it is walking the trie, it will see if the segment we are on is a wildcard segment, and parse out the params for the HTTP Request to use.

If a node/route was found in the router, then it is returned back so it's handler can be called. Will return `nil` if the route was not matched.

### HTTP Handlers

A HTTP Handler that you map gets two parameters passed into it:

* A `HTTPResponse` that contains the request data
* A `HTTPResponseBuilder` to easily build an HTTP Response

A HTTP Handler can return back an `error` if something fails during the response, though right now this isn't used. The goal of this was to maybe build a way to handle errors from the handlers wishing that the `net/http` handlers returned an error as well.

Writing back a response in the handler is as simple as the `HTTPResponseBuilder` which provides a way to write the data needed to the HTTP Request.
Though quite a few helper functions were defined and made available to easily write a response back, without having to write the same thing over and over.

Though if you need to then you can use the `HTTPResponseBuilder` to build out and write responses.

### HTTP Responses

The `HTTPResponseBuilder` given to you in the http handlers gives you a fluent way to build your responses out. It provides the following functions:

* `WithStatus` - Defines the status code and status message for the HTTP Response. Defaults to 200, if not called.
* `WithContent` - Defines the content to send back in the response. If the request is found to have gzip as the accept encoding, then will use will gzip will writing the content
* `WithGzip` - Similar to `WithContent` but enforces gzip to be used as the `content-encoding`. `WithContent` calls this if gzip is used.
* `WithHeader` - Adds a header the response.
* `Write` - This actually writes the response back. This MUST be called, and calling this is essentially ending the response.

