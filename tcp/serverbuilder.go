package tcp

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type TCPServerBuilder struct {
    addr *string
    connHandler func(context.Context, net.Conn)
}

func ServerBuilder() *TCPServerBuilder {
    return &TCPServerBuilder{addr: nil, connHandler: defaultHandler}
}

func (builder *TCPServerBuilder) ListensOn(addr string) {
    builder.addr = &addr
}

func (builder *TCPServerBuilder) UsesHandler(handler func(context.Context, net.Conn)) {
    builder.connHandler = func(ctx context.Context, c net.Conn) {
        defer c.Close()
        handler(ctx, c)
    }
}

// Run will actually start and run the TCPServer and have it start listening to connections.
// This panics if the server can't start listening.
// returns an error if a TCP Connection fails.
func (builder *TCPServerBuilder) Run() error {
    fmt.Println("Starting server...")
    server := TCPServer{}
    var err error

    server.Listener, err = net.Listen("tcp", builder.address())
    if err != nil {
        panic(err)
    }
    server.HandlerFunc = builder.connHandler
    
    server.RegisterCloseSignal()
    fmt.Printf("Listening on %s\n", builder.address())
    err = server.Listen()
    server.Wait(5 * time.Second)

    return err
}


func (builder *TCPServerBuilder) address() string {
    if builder.addr == nil {
        return ""
    }

    return *builder.addr
}
func defaultHandler(ctx context.Context, conn net.Conn) {
    defer conn.Close()
    defer fmt.Println("Connection closed...")

    fmt.Println("Handling Connection...")
    io.Copy(conn, conn)
}
