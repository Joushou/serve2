package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/joushou/serve2"
)

type HTTPHandler struct{}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" || r.Method == "HEAD" {
		return
	}

	fmt.Fprintf(w, "<!DOCTYPE html><html><head></head><body>Welcome to %s</body></html>", r.URL.Path)

}

func main() {

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := serve2.New()

	// Requires host_key (a private key generated by ssh-keygen), and
	// authorized_keys.  Do note that it will log the users in authoized_keys in
	// as the server user.
	ssh := serve2.NewSSHProtoHandler("./host_key", "./authorized_keys")

	// See the HTTPHandler above
	http := serve2.NewHTTPProtoHandler(&HTTPHandler{})

	// Requires cert.pem and key.pem, generated by openssl or
	// http://golang.org/src/crypto/tls/generate_cert.go
	// The protocols listed are used for negatiation, http being required for
	// HTTP over TLS to work.
	tls, err := serve2.NewTLSProtoHandler([]string{"http/1.1"}, "cert.pem", "key.pem")
	if err != nil {
		panic(err)
	}

	// These two are silly, and requires that you write "ECHO" or "DISCARD" when
	// the connection is opened to recognize the protocol, as neither of these
	// actually have any initial request or handshake.
	echo := serve2.NewEchoProtoHandler()
	discard := serve2.NewDiscardProtoHandler()

	server.AddHandlers(ssh, http, tls, echo, discard)
	server.Serve(l)
}
