package main

import (
	"crypto/rand"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {

	// Load config file
	c := LoadConfig()

	h3addr := c.H3Addr
	servername := c.ServerName
	cert := c.CertPath
	key := c.KeyPath
	destination := c.Destinations

	log.Println("server listening on " + h3addr)

	// Generate TLS config for HTTP/3 server
	tconf := tls.Config{Rand: rand.Reader, ServerName: servername, NextProtos: []string{"h3", "h2", "http/1.1"}}
	// Generate QUIC config
	qconf := quic.Config{
		HandshakeIdleTimeout:           time.Duration(c.QUIC.HandshakeIdleTimeout) * time.Second,
		MaxIdleTimeout:                 time.Duration(c.QUIC.MaxIdleTimeout) * time.Second,
		InitialStreamReceiveWindow:     c.QUIC.InitialStreamReceiveWindow,
		MaxStreamReceiveWindow:         c.QUIC.MaxStreamReceiveWindow,
		InitialConnectionReceiveWindow: c.QUIC.InitialConnectionReceiveWindow,
		MaxConnectionReceiveWindow:     c.QUIC.MaxConnectionReceiveWindow,
		MaxIncomingStreams:             c.QUIC.MaxIncomingStreams,
		MaxIncomingUniStreams:          c.QUIC.MaxIncomingUniStreams,
		DisablePathMTUDiscovery:        c.QUIC.DisablePathMTUDiscovery,
		Allow0RTT:                      c.QUIC.Allow0RTT,
	}

	// HTTP/3 Server
	// "QuicConfig: nil" refers to the default configuration for QUIC
	// Handler refers to incoming HTTP request handler
	server := http3.Server{
		Addr:       h3addr,
		QUICConfig: &qconf,
		TLSConfig:  &tconf,
		Handler:    H3Handler(h3addr, destination),
	}

	defer server.Close()

	// Start Listening
	log.Fatalln(server.ListenAndServeTLS(cert, key))
}

// Handle HTTP Request
func H3Handler(H3Addr string, destinations []Destination) http.Handler {
	mux := http.NewServeMux()

	// Register multiple paths
	for _, dest := range destinations {
		mux.HandleFunc(dest.Path, func(w http.ResponseWriter, r *http.Request) {
			// Create an HTTP client on TCP socket
			// {InsecureSkipVerify: true} is required if H1 server Scheme is HTTPS and using self signed certificate
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			h1Client := &http.Client{Transport: tr}

			// Generate HTTP/1.1 request
			h1req := http.Request{Method: r.Method, URL: &url.URL{Scheme: dest.Scheme, Host: dest.Addr, Path: r.URL.Path}}
			// Set H3 request header to h1 request header
			h1req.Header = r.Header
			// Set H1 Extra Headers from config
			MergeMap(h1req.Header, dest.H1ReqHeaders)
			// Set H3 request body to h1 request body
			h1req.Body = r.Body

			// Make HTTP/1.1 request
			response, h1_err := h1Client.Do(&h1req)
			if h1_err != nil {
				log.Println(h1_err.Error())
				w.WriteHeader(500)
				return
			}

			// Set H3 Response Headers
			h3Headers := w.Header()
			MergeMap(h3Headers, response.Header)
			// Set H3 Extra Headers from config
			MergeMap(h3Headers, dest.H3RespHeaders)
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
	}

	return mux
}
