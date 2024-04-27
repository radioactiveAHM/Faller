# Faller

HTTP/3 to HTTP/1.1 proxy server

## Build

To build run `go build`

## Config File

* H3Addr: Refers to the address and port of the HTTP/3 server.

* H1Addr: Indicates the address of the HTTP server where HTTP requests are proxied to.

* ServerName: Your domain address that will be used during the QUIC handshake.

* CertPath: Specifies the path to the certificate file for the HTTP/3 server.

* KeyPath: Specifies the path to the key file for the HTTP/3 server.
