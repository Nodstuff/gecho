package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if fileExists("./ssl/certs/server.crt") {
		startSecureAndInsecure()
	} else {
		startInsecureOnly()
	}
}

func handleEcho() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(buildResponseBody(req))
		if err != nil {
			log.Println(err)
		}
		if _, err = io.Copy(rw, bytes.NewReader(data)); err != nil {
			log.Println(err)
		}
	})
}

func buildResponseBody(req *http.Request) map[string]any {
	rbm := make(map[string]any)
	rb, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(rb, &rbm)
	if err != nil {
		log.Println(err)
	}
	b := make(map[string]any)
	b["statusBody"] = "Healthy"
	b["statusReason"] = fmt.Sprintf("Incoming request was on port %s", getPort(req.Host, req.TLS))
	b["hostname"] = req.Host
	b["uri"] = buildURIResponse(req)
	b["network"] = buildNetworkResponse(req)
	b["ssl"] = buildSSLResponse(req)
	b["requestHeaders"] = buildRequestHeadersResponse(req)
	b["session"] = buildSessionResponse(req)
	b["body"] = rbm
	b["statusCode"] = 200
	return b
}

func buildRequestHeadersResponse(req *http.Request) map[string]any {
	h := make(map[string]any)
	for key, values := range req.Header {
		for _, value := range values {
			h[key] = value
		}
	}
	return h
}

func buildURIResponse(req *http.Request) map[string]any {
	u := make(map[string]any)
	u["httpVersion"] = req.Proto
	u["method"] = req.Method
	u["scheme"] = getScheme(req.TLS)
	u["fullPath"] = req.URL.Path
	u["queryString"] = req.URL.Query().Encode()
	return u
}

func buildNetworkResponse(req *http.Request) map[string]any {
	n := make(map[string]any)
	n["clientPort"] = getPort(req.RemoteAddr, nil)
	n["serverPort"] = getPort(req.Host, req.TLS)
	n["serverAddress"] = req.Host
	n["clientAddress"] = req.RemoteAddr
	return n
}

func buildSessionResponse(req *http.Request) map[string]any {
	s := make(map[string]any)
	s["cookie"] = req.Cookies()
	return s
}

func buildSSLResponse(req *http.Request) map[string]any {
	s := make(map[string]any)
	if req.TLS != nil {
		s["negotiatedProtocol"] = req.TLS.NegotiatedProtocol
		s["cipherSuite"] = req.TLS.CipherSuite
		s["serverName"] = req.TLS.ServerName
		s["version"] = req.TLS.Version
	}
	return s
}

func getScheme(s *tls.ConnectionState) string {
	if s == nil {
		return "http"
	}
	return "https"
}

func getPort(a string, t *tls.ConnectionState) string {
	if t != nil {
		return "443"
	}
	_, p, _ := net.SplitHostPort(a)
	if p == "" {
		p = "80"
	}
	return p
}

func startInsecureOnly() {
	if err := http.ListenAndServe(":8080", handleEcho()); err != nil {
		log.Fatal(err)
	}
}

func startSecureAndInsecure() {
	go func() {
		if err := http.ListenAndServe(":8080", handleEcho()); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServeTLS(":8443", "./ssl/certs/server.crt", "./ssl/certs/server.key", handleEcho()); err != nil {
		log.Fatal(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
