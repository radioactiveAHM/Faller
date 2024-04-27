package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/quic-go/quic-go/http3"
)

func main() {

	// Load config file
	c := LoadConfig()

	h3addr := c.H3Addr
	h1addr := c.H1Addr
	servername := c.ServerName
	cert := c.CertPath
	key := c.KeyPath
	scheme := c.Scheme

	fmt.Println("server listening on " + h3addr)

	// Generate TLS config for HTTP/3 server
	tconf := tls.Config{Rand: rand.Reader, ServerName: servername, NextProtos: []string{"h3", "h2", "http/1.1"}}

	// HTTP/3 Server
	// "QuicConfig: nil" refers to the default configuration for QUIC
	// Handler refers to incoming HTTP request handler
	server := http3.Server{
		Addr:       h3addr,
		QuicConfig: nil,
		TLSConfig:  &tconf,
		Handler:    H3Handler(h1addr, h3addr, scheme),
	}

	defer server.Close()

	// Start Listening
	h3server_err := server.ListenAndServeTLS(cert, key)
	if h3server_err != nil {
		fmt.Println(h3server_err.Error())
		os.Exit(1)
	}
}

func H3Handler(H1Addr string, H3Addr string, scheme string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		h1Client := &http.Client{Transport: tr}
		h1req := http.Request{Method: r.Method, URL: &url.URL{Scheme: scheme, Host: H1Addr, Path: r.URL.Path}, Host: r.Host}
		h1req.Body = r.Body
		response, h1_err := h1Client.Do(&h1req)
		if h1_err != nil {
			fmt.Println(h1_err.Error())
			w.WriteHeader(500)
			return
		}

		// Set HTTP/3 Response
		h3Headers := w.Header()
		for h, v := range response.Header {
			h3Headers.Add(h, strings.Join(v, ";"))
		}

		defer response.Body.Close()
		_, e := io.Copy(w, response.Body)
		if e != nil {
			fmt.Println(e.Error())
			return
		}
	})

	return mux
}
