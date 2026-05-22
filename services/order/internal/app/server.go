package app

import "net/http"

func buildHttpServer(handlers http.Handler) *http.Server {
	return &http.Server{
		Addr:              serverAddr,
		Handler:           handlers,
		ReadHeaderTimeout: timeout,
	}
}
