package minhttp

import (
	"fmt"
	"strings"
)

type RouteNode struct {
    part string
    children map[string]*RouteNode
    handler HttpHandler
    isParam bool
    isRoute bool
}

func newRouteNode() *RouteNode {
    return &RouteNode{
        children: make(map[string]*RouteNode),
    }
}

func routeNodeFromSegment(segment string) *RouteNode {
    return &RouteNode{
        part: segment,
        children: make(map[string]*RouteNode),
        isParam: strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}"),
    }
}

type routePath struct {
    routeNode *RouteNode
    path []string
}

func newRoutePath(node *RouteNode, path []string) *routePath {
    return &routePath{
        routeNode: node,
        path: path,
    }
}

type Router struct {
    getRoot *RouteNode
    postRoot *RouteNode
    putRoot *RouteNode
    patchRoot *RouteNode
}

func NewRouter() *Router {
    return &Router{
        getRoot: newRouteNode(),
        postRoot: newRouteNode(),
        putRoot: newRouteNode(),
        patchRoot: newRouteNode(),
    }
}

func (r *Router) GetRoutes() []string {
    return r.routesFromSubTree(r.getRoot, "/")
}

func (r *Router) PostRoutes() []string {
    return r.routesFromSubTree(r.postRoot, "/")
}

func (r *Router) PutRoutes() []string {
    return r.routesFromSubTree(r.putRoot, "/")
}

func (r *Router) PatchRoutes() []string {
    return r.routesFromSubTree(r.patchRoot, "/")
}

func (r *Router) routesFromSubTree(startingNode *RouteNode, prefix string) []string {
    routes := make([]string, 0, 10)
    seen := make([]*routePath, 0, 10)
    seen = append(seen, newRoutePath(startingNode, make([]string, 0)))

    var node *routePath
    for len(seen) != 0 {
        node, seen = seen[len(seen) - 1], seen[:len(seen) - 1]
        if node.routeNode.isRoute {
            route := fmt.Sprintf("%s%s", prefix, strings.Join(node.path, "/"))
            route = strings.TrimSuffix(route, "/")
            routes = append(routes, route)
        }
        for _, child := range node.routeNode.children {
            newPath := append(node.path, child.part)
            seen = append(seen, newRoutePath(child, newPath))
        }
    }

    return routes
}

func (r *Router) MapGet(template string, handler HttpHandler) {
    r.addRoute(template, r.getRoot, handler)
}

func (r *Router) MapPost(template string, handler HttpHandler) {
    r.addRoute(template, r.postRoot, handler)
}

func (r *Router) MapPut(template string, handler HttpHandler) {
    r.addRoute(template, r.putRoot, handler)
}

func (r *Router) MapPatch(template string, handler HttpHandler) {
    r.addRoute(template, r.patchRoot, handler)
}

func (r *Router) MatchRoute(method, path string) (*RouteNode, map[string]string) {
    segments := parseRoute(path)
    currNode := r.getRoot
    if method == "POST" {
        currNode = r.postRoot
    } else if method == "PUT" {
        currNode = r.putRoot
    } else if method == "PATCH" {
        currNode = r.patchRoot
    }
    params := make(map[string]string)

    for _, segment := range segments {
        if child, exists := currNode.children[segment]; exists {
            currNode = child
            continue
        }

        found := false
        for _, child := range currNode.children {
            if child.isParam {
                paramName := child.part[1:len(child.part) - 1]
                params[paramName] = segment
                currNode = child
                found = true
                break
            }
        }
        if !found {
            return nil, params
        }
    }

    return currNode, params
}

func (r *Router) addRoute(template string, routeRoot *RouteNode, handler HttpHandler) {
    routeSegments := parseRoute(template) 

    currRoute := routeRoot
    for _, segment := range routeSegments {
        if currRoute.children[segment] == nil {
            currRoute.children[segment] = routeNodeFromSegment(segment)
        }

        currRoute = currRoute.children[segment]
    }
    currRoute.isRoute = true
    currRoute.handler = handler
}

func parseRoute(route string) []string {
    route = strings.Trim(route, "/")
    routeComponents := strings.Split(route, "/")

    return routeComponents
}
