package reverseproxy

import (
	"net/http"
	"strings"
)

type Route struct {
	Method         []string
	Path           string
	RewritePath    string
	RequestHeader  http.Header
	ModifyResponse ResponseModifier
}

func NewRoute(methods, path string) Route {
	return Route{
		Method: methodStringToSlice(methods),
		Path:   path,
	}
}

func (r Route) SetRewritePath(path string) Route {
	r.RewritePath = path
	return r
}

func (r Route) SetRequestHeader(header http.Header) Route {
	r.RequestHeader = header
	return r
}

func (r Route) SetModifyResponse(modifier ResponseModifier) Route {
	r.ModifyResponse = modifier
	return r
}

func methodStringToSlice(methods string) []string {
	if methods == "*" {
		return []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}
	}
	return strings.Split(methods, "|")
}
