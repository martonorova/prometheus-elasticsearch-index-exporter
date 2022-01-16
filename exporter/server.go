package exporter

import (
	"fmt"
	"net"
	"net/http"
)

type HttpServerPathHandler struct {
	Path    string
	Handler http.Handler
}

func RunHttpServer(host string, port int, httpHandlers []HttpServerPathHandler) error {
	err := tryOpenPort(host, port)
	if err != nil {
		return listenFailedError(host, port, err)
	}
	for _, httpHandler := range httpHandlers {
		http.Handle(httpHandler.Path, httpHandler.Handler)
	}

	return http.ListenAndServe(fmt.Sprintf("%v:%v", host, port), nil)
}

// Golang's http.ListenAndServe() has an unexpected behaviour when the port is in use:
// Instead of returning an error, it tries to open an IPv6-only listener.
// If this works (because the other application on that port is IPv4-only), no error is returned.
// This is confusing for the user, we want an error if the IPv4 port is in use.
func tryOpenPort(host string, port int) error {
	ln, err := net.Listen("tcp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		return err
	}
	return ln.Close()
}

func listenFailedError(host string, port int, err error) error {
	if len(host) > 0 {
		return fmt.Errorf("cannot bind to %v:%v: %v", host, port, err)
	} else {
		return fmt.Errorf("cannot open port %v: %v", port, err)
	}
}
