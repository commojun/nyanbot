package handler

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/commojun/nyanbot/internal/config"
)

type Server struct {
	APIs []API
	Port int
}

func NewServer(cfg config.Config) (*Server, error) {
	apis, err := newAPIs(cfg)
	return &Server{
		APIs: apis,
		Port: cfg.ServerPort,
	}, err
}

func (server *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	for _, api := range server.APIs {
		newApi := api
		mux.HandleFunc(newApi.MakeHundleFunc())
		log.Printf("Registered API: %s", api.Name)
	}

	port := strconv.Itoa(server.Port)
	httpSrv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("Start server port:%s", port)
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received, shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}
