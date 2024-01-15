package httputil

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMergeRequestHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers []http.Header
		want    http.Header
	}{
		{
			name: "Test with one header",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
		{
			name: "Test with multiple headers",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
				{
					"Accept": []string{"application/json"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
		},
		{
			name: "Test with overlapping headers",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
				{
					"Content-Type": []string{"text/plain"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			MergeRequestHeaders(req, tt.headers...)

			if !reflect.DeepEqual(req.Header, tt.want) {
				t.Errorf("MergeRequestHeaders() = %v, want %v", req.Header, tt.want)
			}
		})
	}
}

func TestMergeResponseHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers []http.Header
		want    http.Header
	}{
		{
			name: "Test with one header",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
		{
			name: "Test with multiple headers",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
				{
					"Accept": []string{"application/json"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
		},
		{
			name: "Test with overlapping headers",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
				{
					"Content-Type": []string{"text/plain"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{Header: http.Header{}}
			MergeResponseHeaders(res, tt.headers...)

			if !reflect.DeepEqual(res.Header, tt.want) {
				t.Errorf("MergeResponseHeaders() = %v, want %v",  res.Header, tt.want)
			}
		})
	}
}

func TestMergeResponseWriterHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers []http.Header
		want    http.Header
	}{
		{
			name: "Test with one header",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
		{
			name: "Test with multiple headers",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
				{
					"Accept": []string{"application/json"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"application/json"},
				"Accept":       []string{"application/json"},
			},
		},
		{
			name: "Test with overlapping headers",
			headers: []http.Header{
				{
					"Content-Type": []string{"application/json"},
				},
				{
					"Content-Type": []string{"text/plain"},
				},
			},
			want: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			MergeResponseWriterHeaders(w, tt.headers...)

			if !reflect.DeepEqual(w.Header(), tt.want) {
				t.Errorf("MergeResponseWriterHeaders() = %v, want %v", w.Header(), tt.want)
			}
		})
	}
}
