package main

import (
	"crypto/rand"
	"crypto/tls"
	"io"
	"log"
	mrand "math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {

	// Load config file
	c := LoadConfig()

	h3addr := c.H3Addr
	domains := c.TLS
	destination := c.Destinations
	trace := c.Trace
	filelog := c.FileLog

	log.Println("server listening on " + h3addr)

	// Generate TLS config for HTTP/3 server
	// Config for multiple Domains
	tconf := tls.Config{
		GetConfigForClient: func(chi *tls.ClientHelloInfo) (*tls.Config, error) {
			if chi.ServerName != "" {
				for _, domain := range domains.Domains {
					if strings.Contains(domain.ServerName, chi.ServerName) {
						// Load Certificate
						cert, cert_e := tls.LoadX509KeyPair(
							domain.CertPath,
							domain.KeyPath,
						)
						if cert_e != nil {
							log.Fatalln(cert_e.Error())
						}
						return &tls.Config{
							Rand:         rand.Reader,
							NextProtos:   []string{"h3", "h2", "http/1.1"},
							Certificates: []tls.Certificate{cert},
						}, nil
					}
				}
			}
			// If no match found return default Certificate
			// No server name support or using ip
			cert, cert_e := tls.LoadX509KeyPair(domains.Default.CertPath, domains.Default.KeyPath)
			if cert_e != nil {
				log.Fatalln(cert_e.Error())
			}
			return &tls.Config{
				Rand:         rand.Reader,
				NextProtos:   []string{"h3", "h2", "http/1.1"},
				Certificates: []tls.Certificate{cert},
			}, nil
		},
	}
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
	// Handler refers to incoming HTTP request handler
	server := http3.Server{
		Addr:       h3addr,
		QUICConfig: &qconf,
		TLSConfig:  &tconf,
		Handler:    H3Handler(h3addr, destination, trace),
	}

	// Logging to file setup
	if filelog.Enable {
		LogFile(&server, filelog.Level)
	}

	defer server.Close()

	// Start Listening
	log.Fatalln(server.ListenAndServe())
}

// Handle HTTP Request
func H3Handler(H3Addr string, destinations []Destination, trace bool) http.Handler {
	mux := http.NewServeMux()

	// Register multiple paths
	for _, dest := range destinations {
		mux.HandleFunc(dest.Path, func(w http.ResponseWriter, r *http.Request) {
			// ID for Tracing
			id := mrand.Intn(65535)

			// Create an HTTP client on TCP socket
			// {InsecureSkipVerify: true} is required if H1 server Scheme is HTTPS and using self signed certificate
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			h1Client := &http.Client{Transport: tr}

			// Trace HTTP/3 incoming request
			if trace {
				TraceLogReq(r.Method, r.URL.Path, r.Proto, r.RemoteAddr, id)
			}

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

			// Trace HTTP/1.1 incoming response
			if trace {
				TraceLogResp(response.Proto, response.StatusCode, id)
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
