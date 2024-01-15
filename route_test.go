package reverseproxy

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewRoute(t *testing.T) {
	type args struct {
		methods string
		path    string
	}
	tests := []struct {
		name string
		args args
		want Route
	}{
		{
			name: "Test with one method",
			args: args{
				methods: "GET",
				path:    "/test",
			},
			want: Route{
				Method: []string{"GET"},
				Path:   "/test",
			},
		},
		{
			name: "Test with multiple methods",
			args: args{
				methods: "GET|POST",
				path:    "/test",
			},
			want: Route{
				Method: []string{"GET", "POST"},
				Path:   "/test",
			},
		},
		{
			name: "Test with all methods",
			args: args{
				methods: "*",
				path:    "/test",
			},
			want: Route{
				Method: []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"},
				Path:   "/test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoute(tt.args.methods, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_SetRewritePath(t *testing.T) {
	type fields struct {
		RewritePath string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Route
	}{
		{
			name: "Test with empty path",
			fields: fields{
				RewritePath: "",
			},
			args: args{
				path: "",
			},
			want: &Route{
				RewritePath: "",
			},
		},
		{
			name: "Test with non-empty path",
			fields: fields{
				RewritePath: "",
			},
			args: args{
				path: "/new-path",
			},
			want: &Route{
				RewritePath: "/new-path",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				RewritePath: tt.fields.RewritePath,
			}
			if got := r.SetRewritePath(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.SetRewritePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_SetRequestHeader(t *testing.T) {
	type fields struct {
		RequestHeader http.Header
	}
	type args struct {
		header http.Header
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Route
	}{
		{
			name: "Test with empty header",
			fields: fields{
				RequestHeader: http.Header{},
			},
			args: args{
				header: http.Header{},
			},
			want: &Route{
				RequestHeader: http.Header{},
			},
		},
		{
			name: "Test with one header",
			fields: fields{
				RequestHeader: http.Header{},
			},
			args: args{
				header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			want: &Route{
				RequestHeader: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
		},
		{
			name: "Test with multiple headers",
			fields: fields{
				RequestHeader: http.Header{},
			},
			args: args{
				header: http.Header{
					"Content-Type": []string{"application/json"},
					"Accept":       []string{"application/json"},
				},
			},
			want: &Route{
				RequestHeader: http.Header{
					"Content-Type": []string{"application/json"},
					"Accept":       []string{"application/json"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				RequestHeader: tt.fields.RequestHeader,
			}
			if got := r.SetRequestHeader(tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Route.SetRequestHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoute_SetModifyResponse(t *testing.T) {
	mockModifier := func(resp *http.Response) error {
		resp.Header.Set("X-Test", "modified")
		return nil
	}
	tests := []struct {
		name             string
		modifyResponseIn ResponseModifier 
		expectedHeader   string       
	}{
		{
			name:             "Test with nil modifier",
			modifyResponseIn: nil,
			expectedHeader:   "",
		},
		{
			name:             "Test with non-nil modifier",
			modifyResponseIn: mockModifier,
			expectedHeader:   "modified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{}
			r.SetModifyResponse(tt.modifyResponseIn)
			mockResp := &http.Response{Header: http.Header{}}
			if r.ModifyResponse != nil {
				err := r.ModifyResponse(mockResp)
				if err != nil {
					t.Errorf("Modifier returned an error: %v", err)
				}
			}
			if gotHeader := mockResp.Header.Get("X-Test"); gotHeader != tt.expectedHeader {
				t.Errorf("Route.SetModifyResponse() not working correctly")
			}
		})
	}
}
