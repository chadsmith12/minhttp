package minhttp

import (
	"context"
	"net"

	"github.com/chadsmith12/minhttp/tcp"
)

var defaultRouter *Router = NewRouter()

type HttpHandler = func(HttpRequest, ResponseWriter) error

type Server struct {
    Addr string
    router *Router
}

func (server *Server) ListenAndServe() error {
    tcpBuilder := tcp.ServerBuilder()
    tcpBuilder.ListensOn(server.Addr)
    tcpBuilder.UsesHandler(server.defaultHttpHanlder)

    return tcpBuilder.Run()
}

func ListenAndServe(addr string) error {
    server := &Server{ Addr: addr, router: defaultRouter }

    return server.ListenAndServe()
}

func MapGet(template string, handler HttpHandler) {
    defaultRouter.MapGet(template, handler)
}

func MapPost(template string, handler HttpHandler) {
    defaultRouter.MapPost(template, handler)
}

func MapPut(template string, handler HttpHandler) {
    defaultRouter.MapPut(template, handler)
}

func MapPatch(template string, handler HttpHandler) {
    defaultRouter.MapPatch(template, handler)
}

func (s *Server) defaultHttpHanlder(ctx context.Context, conn net.Conn) {
    req, err := ReadRequest(conn)
    resp := &httpResponse{ request: &req, headers: NewHeadersCollection(), conn: conn }
    if err != nil {
	WriteBadRequest(resp)
    }
    route, params := s.router.MatchRoute(req.Method, req.Path)
    if route == nil {
        WriteNotFound(resp)
	return
    }
    req.Params = params
    route.handler(req, resp)
}
