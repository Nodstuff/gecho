package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8080", handleEcho()); err != nil {
		log.Fatal("He's dead, Jim!")
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
	b["method"] = req.Method
	b["url"] = req.URL.String()
	b["proto"] = req.Proto
	b["content-length"] = req.ContentLength
	b["host"] = req.Host
	b["remote-addr"] = req.RemoteAddr
	b["headers"] = req.Header
	b["params"] = req.URL.Query()
	b["body"] = rbm
	return b
}
