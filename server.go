package minhttp

import (
	"context"
	"net"
	"strings"

	"github.com/chadsmith12/minhttp/tcp"
)

type Server struct {
    Addr string
}

func (server *Server) ListenAndServe() error {
    tcpBuilder := tcp.ServerBuilder()
    tcpBuilder.ListensOn(server.Addr)
    tcpBuilder.UsesHandler(defaultHttpHanlder)

    return tcpBuilder.Run()
}

func ListenAndServe(addr string) error {
    server := &Server{ Addr: addr }

    return server.ListenAndServe()
}

func defaultHttpHanlder(ctx context.Context, conn net.Conn) {
    req, err := ReadRequest(conn)

    if err != nil {
        WriteBadRequest(conn)
    }
    if req.Path == "/" {
        WriteOk(conn)
    } else if strings.HasPrefix(req.Path, "/echo") {
        components := strings.Split(req.Path, "/echo/")
        WriteText(conn, components[1]) 
    } else if strings.HasPrefix(req.Path, "/user-agent") {
        userAgent := req.Headers.UserAgent()
        WriteText(conn, userAgent)
    } else {
        WriteNotFound(conn)
    }
}
