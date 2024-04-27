package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	fmt.Println("Server Started")

	c := LoadConfig()

	h3addr := c.H3Addr
	h1addr := c.H1Addr
	servername := c.ServerName
	cert := c.CertPath
	key := c.KeyPath

	tconf := tls.Config{Rand: rand.Reader, ServerName: servername, NextProtos: []string{"h3", "h2", "http/1.1"}}

	server := http3.Server{
		Addr:       h3addr,
		QuicConfig: nil,
		TLSConfig:  &tconf,
		Handler:    H3Handler(h1addr, h3addr),
	}

	defer server.Close()

	h3server_err := server.ListenAndServeTLS(cert, key)
	if h3server_err != nil {
		fmt.Println(h3server_err.Error())
		os.Exit(1)
	}
}

func H3Handler(H1Addr string, H3Addr string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		h1Client := http.DefaultClient
		h1req := http.Request{Method: r.Method, URL: &url.URL{Scheme: "http", Host: H1Addr, Path: r.URL.Path}, Host: r.Host}
		h1req.Body = r.Body
		respone, h1_err := h1Client.Do(&h1req)
		if h1_err != nil {
			fmt.Println(h1_err.Error())
			w.WriteHeader(500)
			return
		}

		h3headers := http.Header{}
		for h, v := range respone.Header {
			h3headers.Add(h, strings.Join(v, " "))
		}
		h3respone := http.Response{Proto: "HTTP/3", StatusCode: respone.StatusCode, Header: h3headers, ContentLength: respone.ContentLength}
		h3respone.Body = respone.Body
		e := h3respone.Write(w)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
	})

	return mux
}
