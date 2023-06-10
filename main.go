package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if fileExists("./ssl/certs/server.crt") {
		startSecureAndInsecure(ctx)
	} else {
		startInsecureOnly()
	}
}

func echoHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		d, err := json.Marshal(buildResponseBody(req))
		if err != nil {
			log.Println(err)
		}

		handleResponseHeaders(rw, req)

		rw.WriteHeader(checkRequestedStatusHeader(req.Method, req.Header))

		if _, err = io.Copy(rw, bytes.NewReader(d)); err != nil {
			log.Println(err)
		}
	})
}

func handleResponseHeaders(rw http.ResponseWriter, req *http.Request) {
	req.Header.Del("Content-Length")
	rw.Header().Set("Content-Type", "application/json")
	for key, values := range req.Header {
		for _, value := range values {
			rw.Header().Set(key, value)
		}
	}
}

func buildResponseBody(r *http.Request) map[string]any {
	rbm := make(map[string]any)
	if r.ContentLength > 0 {
		rb, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		err = json.Unmarshal(rb, &rbm)
		if err != nil {
			log.Println(err)
		}
	}
	b := make(map[string]any)
	b["statusBody"] = "Healthy"
	b["statusReason"] = fmt.Sprintf("Incoming request was on port %s", getPort(r.Host, r.TLS))
	b["hostname"] = r.Host
	b["uri"] = buildURI(r)
	b["network"] = buildNetwork(r)
	b["ssl"] = buildSSL(r)
	b["requestHeaders"] = buildRequestHeaders(r)
	b["session"] = buildSession(r)
	b["body"] = rbm
	b["statusCode"] = http.StatusOK
	return b
}

func buildRequestHeaders(r *http.Request) map[string]any {
	h := make(map[string]any)
	for k, vs := range r.Header {
		for _, v := range vs {
			h[k] = v
		}
	}
	return h
}

func buildURI(r *http.Request) map[string]any {
	u := make(map[string]any)
	u["httpVersion"] = r.Proto
	u["method"] = r.Method
	u["scheme"] = getScheme(r.TLS)
	u["fullPath"] = r.URL.Path
	u["queryString"] = r.URL.Query().Encode()
	return u
}

func buildNetwork(r *http.Request) map[string]any {
	n := make(map[string]any)
	n["clientPort"] = getPort(r.RemoteAddr, nil)
	n["serverPort"] = getPort(r.Host, r.TLS)
	n["serverAddress"] = r.Host
	n["clientAddress"] = r.RemoteAddr
	return n
}

func buildSession(r *http.Request) map[string]any {
	s := make(map[string]any)
	s["cookie"] = r.Cookies()
	return s
}

func buildSSL(r *http.Request) map[string]any {
	s := make(map[string]any)
	if r.TLS != nil {
		s["negotiatedProtocol"] = r.TLS.NegotiatedProtocol
		s["cipherSuite"] = r.TLS.CipherSuite
		s["serverName"] = r.TLS.ServerName
		s["version"] = r.TLS.Version
	}
	return s
}

func getScheme(s *tls.ConnectionState) string {
	if s == nil {
		return "http"
	}
	return "https"
}

func getPort(a string, s *tls.ConnectionState) string {
	_, p, _ := net.SplitHostPort(a)
	if p == "" && s != nil {
		return "443"
	} else if p == "" {
		return "80"
	}
	return p
}

func fileExists(f string) bool {
	i, err := os.Stat(f)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !i.IsDir()
}

func checkRequestedStatusHeader(m string, h http.Header) int {
	status, err := strconv.Atoi(h.Get("X-Requested-Status"))
	if err != nil {
		return getDefaultStatus(m)
	}
	return status
}

func getDefaultStatus(m string) int {
	switch m {
	case http.MethodPost:
		return http.StatusCreated
	case http.MethodDelete:
		return http.StatusNoContent
	default:
		return http.StatusOK
	}
}

func startInsecureOnly() {
	if err := http.ListenAndServe(":8080", echoHandler()); err != nil {
		log.Fatal(err)
	}
}

func startSecureAndInsecure(ctx context.Context) {
	h := http.Server{
		Addr:    ":8080",
		Handler: echoHandler(),
	}

	go func() {
		<-ctx.Done()
		err := h.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServeTLS(":8443", "./ssl/certs/server.crt", "./ssl/certs/server.key", echoHandler()); err != nil {
		log.Fatal(err)
	}
}
