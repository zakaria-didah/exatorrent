package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zakaria-didah/exatorrent/internal/core"
	"github.com/zakaria-didah/exatorrent/internal/web"
)

func main() {
	core.Initialize()

	mux := http.NewServeMux()
	mux.Handle("/api/auth", core.RateLimit(core.CORSMiddleware(http.HandlerFunc(core.AuthCheck))))
	mux.Handle("/api/auth/logout", core.CORSMiddleware(http.HandlerFunc(core.LogoutHandler)))
	mux.Handle("/api/socket", core.CORSMiddleware(http.HandlerFunc(core.SocketAPI)))
	mux.Handle("/api/stream/", core.CORSMiddleware(http.HandlerFunc(core.StreamFile)))
	mux.Handle("/api/torrent/", core.CORSMiddleware(http.HandlerFunc(core.TorrentServe)))
	mux.Handle("/api/poster", core.CORSMiddleware(http.HandlerFunc(core.PosterHandler)))
	mux.Handle("/api/signup-request", core.RateLimit(core.CORSMiddleware(http.HandlerFunc(core.SignupRequestHandler))))
	mux.Handle("/api/thumbnail/", core.CORSMiddleware(http.HandlerFunc(core.ThumbnailHandler)))
	mux.Handle("/", web.FrontEndHandler)

	handler := core.RequestLogger(mux)

	if core.Flagconfig.UnixSocket != "" {
		core.Info.Println("Starting server at Path (Unix Socket)", core.Flagconfig.UnixSocket)
		usock, err := net.Listen("unix", core.Flagconfig.UnixSocket)
		if err != nil {
			core.Err.Fatalln("Failed listening", err)
		}
		go func() {
			if core.Flagconfig.TLSCertPath != "" && core.Flagconfig.TLSKeyPath != "" {
				core.Info.Println("Serving the HTTPS with TLS Cert ", core.Flagconfig.TLSCertPath, " and TLS Key", core.Flagconfig.TLSKeyPath)
				core.Err.Fatal(http.ServeTLS(usock, handler, core.Flagconfig.TLSCertPath, core.Flagconfig.TLSKeyPath))
			} else {
				core.Err.Fatal(http.Serve(usock, handler))
			}
		}()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		_ = usock.Close()
	} else {
		core.Info.Println("Starting server on", core.Flagconfig.ListenAddress)
		if core.Flagconfig.TLSCertPath != "" && core.Flagconfig.TLSKeyPath != "" {
			core.Info.Println("Serving the HTTPS with TLS Cert ", core.Flagconfig.TLSCertPath, " and TLS Key", core.Flagconfig.TLSKeyPath)
			core.Err.Fatal(http.ListenAndServeTLS(core.Flagconfig.ListenAddress, core.Flagconfig.TLSCertPath, core.Flagconfig.TLSKeyPath, handler))
		} else {
			core.Err.Fatal(http.ListenAndServe(core.Flagconfig.ListenAddress, handler))
		}
	}
}
