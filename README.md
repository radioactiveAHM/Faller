# Faller

HTTP/3 to HTTP/1.1 proxy server

## Roadmap

- [ ] Multiple path handler

## Build

To build run `go build`

## Config File

* H3Addr: Refers to the address and port of the HTTP/3 server.

* H1Addr: Indicates the address of the HTTP server where HTTP requests are proxied to.

* ServerName: Your domain address that will be used during the QUIC handshake.

* CertPath: Specifies the path to the certificate file for the HTTP/3 server.

* KeyPath: Specifies the path to the key file for the HTTP/3 server.

* Scheme: Specifies the Scheme of the HTTP/1.1 server.

> [!WARNING]
> To inform the browser that your web server supports HTTP/3, you should include the following Alt-Svc HTTP header in your server’s response: `Alt-Svc: h3=":443", h3-29=":443"` This header indicates that HTTP/3 is available on UDP port 443 at the same host name that was used to retrieve this response. By including this header, you allow clients to establish QUIC connections to that destination and continue communicating with the origin using HTTP/3 instead of the initial HTTP version. [Link-to-article](https://www.ietf.org/archive/id/draft-duke-httpbis-quic-version-alt-svc-03.html)
