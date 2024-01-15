package httputil

import "net/http"

func MergeRequestHeaders(r *http.Request, headers ...http.Header) {
	for _, header := range headers {
		for k, v := range header {
			r.Header[k] = v
		}
	}
}

func MergeResponseHeaders(r *http.Response, headers ...http.Header) {
	for _, header := range headers {
		for k, v := range header {
			r.Header[k] = v
		}
	}
}

func MergeResponseWriterHeaders(w http.ResponseWriter, headers ...http.Header) {
	for _, header := range headers {
		for k, lines := range header {
			for i, line := range lines {
				if i == 0 {
					w.Header().Set(k, line)
					continue
				}
				w.Header().Add(k, line)
			}
		}
	}
}
