<img src="./faller_ico.webp" width="300">

# Faller

Faller is an HTTP/3 proxy server that allows you to proxy HTTP/3 requests to other web applications. Similar to NGINX’s reverse proxy functionality, Faller supports multipath configuration and manual HTTP header settings. Let’s break down its key features:

1. HTTP/3 Proxy: Faller acts as an intermediary between clients and web servers, forwarding HTTP/3 requests to the appropriate destination. It ensures secure and efficient communication.
2. Multipath Configuration: Faller allows you to configure multiple paths for routing requests. This flexibility enables load balancing and efficient distribution of traffic across different endpoints.
3. Manual HTTP Header Setting: You can customize HTTP headers within Faller. This feature is useful for scenarios where specific headers need modification or addition.

## Roadmap

- [x] Multiple Path Handler
- [x] Enabling Custom HTTP Headers for HTTP/3 Responses and HTTP/1 Requests
- [x] Configuration options for QUIC in config file
- [x] Tracing option
- [x] Logging to file option
- [x] Multi Domain

## Build

To build run `go build` in project directory

To build stripped run `go build -ldflags "-w -s"` in project directory

## Config File

- H3Addr: Refers to the address and port of the HTTP/3 server.

- TLS: Tls configuration

  - Default: Certificate for clients without SNI or using IP
    - CertPath: Specifies the path to the certificate.
    - KeyPath: Specifies the path to the key.

  - Domains (optional): List of server names with specified certificate.
    - ServerName: Ensure that the server name matches the client SNI (Server Name Indication) to return the specified certificate. If all matches failed the Default certificate is returned.
    - CertPath: Specifies the path to the certificate.
    - KeyPath: Specifies the path to the key.
    - SubDomainsSupport: This feature indicates whether the certificate covers all subdomains.

- Destinations: Configuring HTTP Path and HTTP/1 Server: Setting Up Multiple Paths to Proxy Requests to Different HTTP/1 Servers

  - Addr: Indicates the address of the HTTP server where HTTP requests are proxied to.
  - Scheme: Specifies the Scheme of the HTTP/1.1 server.
  - Path: If the HTTP path matches this specified path, the HTTP request will be proxied to the designated address.
  - H3RespHeaders: Extra Header (Set as H3 Response Header)
  - H1ReqHeaders: Extra Header (Set as HTTP/1.1 Request Header)

- Trace: If enabled, this utility will print the incoming HTTP/3 request line and the server’s HTTP/1.1 response status line to the standard output (stdout).

- FileLog: Log HTTP/3 event statuses to a file if enabled. Valid log level variants include: Debug, Error, Info, and Warn.

> [!WARNING]
> To inform the browser that your web server supports HTTP/3, you should include the following Alt-Svc HTTP header in your server’s response: `Alt-Svc: h3=":443", h3-29=":443"` This header indicates that HTTP/3 is available on UDP port 443 at the same host name that was used to retrieve this response. By including this header, you allow clients to establish QUIC connections to that destination and continue communicating with the origin using HTTP/3 instead of the initial HTTP version. [Link-to-article](https://www.ietf.org/archive/id/draft-duke-httpbis-quic-version-alt-svc-03.html)
