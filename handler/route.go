package handler

import (
	"net/http"
	"regexp"
)

type ctxKey struct{}

var routes []route

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func newRouter(method string, pattern string, handler http.HandlerFunc) route {
	return route{
		method:  method,
		regex:   regexp.MustCompile("^" + pattern + "$"),
		handler: handler,
	}
}

func addRoutes(r []route) {
	routes = append(routes, r...)
}
