# Go Reverse Proxy

This project is a lightweight and efficient reverse proxy module written in Go. Besides directly forwarding incoming requests to the specified remote,
it allows for rewriting of the request path and defining custom headers.

## Features

- Easy configuration via a simple Go API.
- Fine-grained control over which paths are passed through to the backend.
- Support for rewriting of the request path.
- Customizable request and response headers.
- Integrated health check and load measurement functionality.

## Installation

To install `go-reverse-proxy`, use the following command:

	go get -u github.com/secondtruth/go-reverse-proxy

## Usage

```go
package main

import (
	"crypto/tls"
	"log"
	"net/http"

	reverseproxy "github.com/secondtruth/go-reverse-proxy"
)

func main() {
	rp, err := reverseproxy.New("http://my-backend:8000")
	if err != nil {
		log.Fatal(err)
	}
	rp.RequestHeader = http.Header{
		"Authorization": []string{"Bearer abc123"},
	}
	rp.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	rp.PassPath("*", "/")
	rp.PassPaths("HEAD|GET|POST", "/api/version", "/api/posts")
	rp.RewritePath("HEAD|GET|POST", "/posts", "/api/posts")

	log.Fatal(http.ListenAndServe(":8080", rp))
}
```
