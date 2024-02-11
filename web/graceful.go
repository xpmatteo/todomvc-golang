package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// GracefulListenAndServe will ensure that when the server receives a SIGTERM signal,
// it will shut down gracefully.  It will wait for all connections, both open and idle
// to be closed before shutting down the process.
//
// There is no timeout enforced, because it is the job of the container to do that.
// For instance, Kubernetes will eventually forcefully kill a pod after waiting
// for a configured timeout for it to exit cleanly.
//
// The implementation is taken from [https://pkg.go.dev/net/http#Server.Shutdown]
func GracefulListenAndServe(addr string, handler http.Handler) {
	server := &http.Server{Addr: addr, Handler: handler}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener
		log.Fatal(err)
	}

	<-idleConnsClosed
}
