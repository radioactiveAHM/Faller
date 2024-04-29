package main

import (
	"crypto/rand"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
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

	log.Println("server listening on " + h3addr)

	// Generate TLS config for HTTP/3 server
	tconf := tls.Config{Rand: rand.Reader, ServerName: servername, NextProtos: []string{"h3", "h2", "http/1.1"}}

	// HTTP/3 Server
	// "QuicConfig: nil" refers to the default configuration for QUIC
	// Handler refers to incoming HTTP request handler
	server := http3.Server{
		Addr:       h3addr,
		QUICConfig: nil,
		TLSConfig:  &tconf,
		Handler:    H3Handler(h1addr, h3addr, scheme),
	}

	defer server.Close()

	// Start Listening
	log.Fatalln(server.ListenAndServeTLS(cert, key))
}

// Handle HTTP Request
func H3Handler(H1Addr string, H3Addr string, scheme string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		// Create an HTTP client on TCP socket
		// {InsecureSkipVerify: true} is required if H1 server Scheme is HTTPS and using self signed certificate
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		h1Client := &http.Client{Transport: tr}

		// Generate HTTP/1.1 request
		h1req := http.Request{Method: r.Method, URL: &url.URL{Scheme: scheme, Host: H1Addr, Path: r.URL.Path}}
		// Set H3 request header to h1 request header
		h1req.Header = r.Header
		// Set H3 request body to h1 request body
		h1req.Body = r.Body

		// Make HTTP/1.1 request
		response, h1_err := h1Client.Do(&h1req)
		if h1_err != nil {
			log.Println(h1_err.Error())
			w.WriteHeader(500)
			return
		}

		// Set HTTP/3 Response
		h3Headers := w.Header()
		for h, v := range response.Header {
			h3Headers.Add(h, strings.Join(v, ","))
		}
		defer response.Body.Close()
		// Write H1 Response StatusCode to H3  Response StatusCode
		w.WriteHeader(response.StatusCode)
		// Write H1 Response Body to H3 Response Body
		_, e := io.Copy(w, response.Body)
		if e != nil {
			log.Println(e.Error())
			return
		}
	})

	return mux
}
