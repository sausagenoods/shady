package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func initRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(100 * time.Second))
	r.Use(middleware.Logger)
	r.Post("/encrypt", encryptHandler)
	r.Get("/decrypt/{id}", decryptHandler)
	r.Post("/callback/{id}", moneroPayCallbackHandler)
	return r
}

func runServer() {
	h2s := &http2.Server{}
	srv := &http.Server{
		Addr: Conf.bindAddr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 30 * time.Second,
		Handler: h2c.NewHandler(initRouter(), h2s),
	}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30 * time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out. Forcing exit.")
			}
		}()
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-serverCtx.Done()
}
