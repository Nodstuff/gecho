package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func Test_buildRequestHeaders(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "all headers copied",
			args: args{r: &http.Request{Header: map[string][]string{
				"Host":            {"127.0.0.1"},
				"X-Forwarded-For": {"192.168.1.1"},
				"Authorization":   {"Basic amltOm15c2VjcmV0cGFzc3dvcmQ="},
			}}},
			want: map[string]any{
				"Host":            "127.0.0.1",
				"X-Forwarded-For": "192.168.1.1",
				"Authorization":   "Basic amltOm15c2VjcmV0cGFzc3dvcmQ=",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildRequestHeaders(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildRequestHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildSSL(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "all values present",
			args: args{r: &http.Request{TLS: &tls.ConnectionState{
				NegotiatedProtocol: "https",
				CipherSuite:        1234,
				ServerName:         "gecho",
				Version:            99,
			}}},
			want: map[string]any{
				"negotiatedProtocol": "https",
				"cipherSuite":        uint16(1234),
				"serverName":         "gecho",
				"version":            uint16(99),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildSSL(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				fmt.Println(got)
				fmt.Println(tt.want)
				t.Errorf("buildSSL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildSession(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "all session values copied",
			args: args{r: &http.Request{Header: map[string][]string{
				"Cookie": {"some-cookie-value"},
			}}},
			want: map[string]any{
				"cookie": []*http.Cookie{{Name: "some-cookie-value"}},
			},
		},
		{
			name: "empty session values",
			args: args{r: &http.Request{}},
			want: map[string]any{
				"cookie": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildSession(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildURI(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "no TLS or query params",
			args: args{r: &http.Request{
				Proto:  "HTTP/1.1",
				Method: http.MethodGet,
				URL: &url.URL{
					Path: "/path/to/thing",
				},
			}},
			want: map[string]any{
				"fullPath":    "/path/to/thing",
				"httpVersion": "HTTP/1.1",
				"method":      http.MethodGet,
				"queryString": "",
				"scheme":      "http",
			},
		},
		{
			name: "with TLS and query params",
			args: args{r: &http.Request{
				TLS:    &tls.ConnectionState{},
				Proto:  "HTTP/1.1",
				Method: http.MethodGet,
				URL: &url.URL{
					Path:     "/path/to/thing?params=true",
					RawQuery: "params=true",
				},
			}},
			want: map[string]any{
				"fullPath":    "/path/to/thing?params=true",
				"httpVersion": "HTTP/1.1",
				"method":      http.MethodGet,
				"queryString": "params=true",
				"scheme":      "https",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildURI(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileExists(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "File does not exist",
			args: args{f: "./madeup-file.txt"},
			want: false,
		},
		{
			name: "File does exist",
			args: args{f: "./main.go"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.args.f); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPort(t *testing.T) {
	type args struct {
		a string
		s *tls.ConnectionState
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "port 443",
			args: args{
				a: "127.0.0.1:443",
				s: &tls.ConnectionState{},
			},
			want: "443",
		},
		{
			name: "port 443",
			args: args{
				a: "127.0.0.1",
				s: &tls.ConnectionState{},
			},
			want: "443",
		},
		{
			name: "port 80",
			args: args{
				a: "127.0.0.1",
				s: nil,
			},
			want: "80",
		},
		{
			name: "port 8443",
			args: args{
				a: "127.0.0.1:8443",
				s: &tls.ConnectionState{},
			},
			want: "8443",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPort(tt.args.a, tt.args.s); got != tt.want {
				t.Errorf("getPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getScheme(t *testing.T) {
	type args struct {
		s *tls.ConnectionState
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "https scheme",
			args: args{s: &tls.ConnectionState{}},
			want: "https",
		},
		{
			name: "http scheme",
			args: args{s: nil},
			want: "http",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getScheme(tt.args.s); got != tt.want {
				t.Errorf("getScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}
