// Package reverseproxy provides a reverse proxy implementation in Go.
// It allows you to create a reverse proxy server that forwards HTTP requests to a remote server.
// The package includes features like response modification, error handling, and request rewriting.
package reverseproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/haoxins/rewrite"
	"github.com/julienschmidt/httprouter"
	"github.com/secondtruth/go-reverse-proxy/health"
	httputilx "github.com/secondtruth/go-reverse-proxy/httputil"
)

// ResponseModifier is a function that modifies the HTTP response.
type ResponseModifier func(*http.Response) error

// ResponseModifierMap is a map of [method][path]ResponseModifier.
type ResponseModifierMap map[string]map[string]ResponseModifier

// HttpErrorHandler is a function that handles errors occurring in HTTP request handlers.
type HttpErrorHandler func(http.ResponseWriter, *http.Request, error)

// ReverseProxyMux is a reverse proxy with a request path multiplexer.
type ReverseProxyMux struct {
	proxy     *httputil.ReverseProxy
	remote    *url.URL
	router    *httprouter.Router
	modifiers ResponseModifierMap
	health    *health.HealthCheck
	load      int32

	Transport               http.RoundTripper
	RequestHeader           http.Header
	ModifyResponse          ResponseModifier
	ErrorHandler            HttpErrorHandler
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
}

// New creates a new ReverseProxyMux with the specified remote URL.
func New(remote string) (*ReverseProxyMux, error) {
	remoteUrl, err := url.Parse(remote)
	if err != nil {
		return nil, err
	}
	pm := &ReverseProxyMux{
		proxy:     httputil.NewSingleHostReverseProxy(remoteUrl),
		remote:    remoteUrl,
		router:    httprouter.New(),
		modifiers: make(ResponseModifierMap),
		health:    health.NewHealthCheck(remoteUrl),
	}
	return pm, nil
}

// ServeHTTP handles the HTTP request.
func (pm *ReverseProxyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&pm.load, 1)
	defer atomic.AddInt32(&pm.load, -1)

	pm.proxy.ModifyResponse = func(r *http.Response) error {
		if pm.ModifyResponse != nil {
			if err := pm.ModifyResponse(r); err != nil {
				return err
			}
		}
		if modifier, ok := pm.modifiers[r.Request.Method][r.Request.URL.Path]; ok {
			if err := modifier(r); err != nil {
				return err
			}
		}
		return nil
	}
	if pm.Transport != nil {
		pm.proxy.Transport = pm.Transport
	}

	pm.router.NotFound = pm.NotFoundHandler
	pm.router.MethodNotAllowed = pm.MethodNotAllowedHandler

	if pm.ErrorHandler != nil {
		pm.proxy.ErrorHandler = pm.ErrorHandler
		pm.router.PanicHandler = func(w http.ResponseWriter, r *http.Request, val any) {
			pm.ErrorHandler(w, r, fmt.Errorf("%v", val))
		}
	}

	pm.router.ServeHTTP(w, r)
}

// HandlePath registers a route.
func (pm *ReverseProxyMux) HandlePath(route Route) *ReverseProxyMux {
	for _, method := range route.Method {
		pm.router.HandlerFunc(method, route.Path, func(w http.ResponseWriter, r *http.Request) {
			r.Header.Set("X-Forwarded-Proto", r.URL.Scheme)
			r.Header.Set("X-Forwarded-Host", r.Host)
			r.Host = pm.remote.Host
			if route.RewritePath != "" {
				rewriter, err := rewrite.NewRule(route.Path, route.RewritePath)
				if err != nil {
					pm.ErrorHandler(w, r, err)
					return
				}
				rewriter.Rewrite(r)
			}
			httputilx.MergeRequestHeaders(r, pm.RequestHeader, route.RequestHeader)

			pm.proxy.ServeHTTP(w, r)
		})
		if route.ModifyResponse != nil {
			if pm.modifiers[method] == nil {
				pm.modifiers[method] = make(map[string]ResponseModifier)
			}
			pm.modifiers[method][route.Path] = route.ModifyResponse
		}
	}
	return pm
}

// PassPath registers a path with the specified HTTP methods.
func (pm *ReverseProxyMux) PassPath(methods, path string) *ReverseProxyMux {
	return pm.HandlePath(NewRoute(methods, path))
}

// PassPaths registers multiple paths with the specified HTTP methods.
func (pm *ReverseProxyMux) PassPaths(methods string, paths ...string) *ReverseProxyMux {
	for _, path := range paths {
		pm.PassPath(methods, path)
	}
	return pm
}

// PassAnyPath registers all possible paths with the specified HTTP methods.
func (pm *ReverseProxyMux) PassAnyPath(methods string) *ReverseProxyMux {
	return pm.PassPath(methods, "/*path")
}

// PassAnyPathUnder registers all possible paths under the specified parent paths with the specified HTTP methods.
func (pm *ReverseProxyMux) PassAnyPathUnder(methods string, paths ...string) *ReverseProxyMux {
	for _, path := range paths {
		pm.PassPath(methods, filepath.Join(path, "/*path"))
	}
	return pm
}

// RewritePath registers a route with the specified HTTP methods and source path, and rewrites the request path to the target path.
func (pm *ReverseProxyMux) RewritePath(methods, sourcePath, targetPath string) *ReverseProxyMux {
	route := NewRoute(methods, sourcePath)
	route.RewritePath = targetPath
	return pm.HandlePath(route)
}

// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
func (p *ReverseProxyMux) IsAvailable() bool {
	return p.health.IsAvailable()
}

// SetHealthCheckFunc sets the passed check func as the algorithm of checking the origin availability
func (p *ReverseProxyMux) SetHealthCheckFunc(check func(addr *url.URL) bool, period time.Duration) {
	p.health.SetCheckFunc(check, period)
}

// GetLoad returns the number of requests being served by the proxy at the moment
func (p *ReverseProxyMux) GetLoad() int32 {
	return atomic.LoadInt32(&p.load)
}
