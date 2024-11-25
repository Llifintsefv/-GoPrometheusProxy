package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	srv *http.Server
	port string
	shutdownTimeout time.Duration
}

func NewServer(handler http.Handler,port string) *Server{
	address := net.JoinHostPort("0.0.0.0",port)
	timeoutString := os.Getenv("SERVER_SHUTDOWN_TIMEOUT")
	timeout := defaultShutdownTimeout
	if timeoutString != "" {
		timeoutInt,err := strconv.Atoi(timeoutString)
		if err != nil {
			timeout = time.Duration(timeoutInt) * time.Second
		} else {
			log.Printf("Invalid SERVER_SHUTDOWN_TIMEOUT: %v, default: %v",err,defaultShutdownTimeout)
		}
	} 
	src := &http.Server{
		Addr: address,
		Handler: handler,
	}

	return &Server{
		srv: src,
		port: port,
		shutdownTimeout: timeout,
	}

}

func (s *Server) Run(ctx context.Context) error{
	log.Printf("Starting server on port %s",s.port)
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err,http.ErrServerClosed) {
			log.Printf("Liste and serve error := %v ",err)
		}
	}()
	<-ctx.Done()

	log.Println("Shutting down server")
	downCtx,cancel := context.WithTimeout(context.Background(),s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(downCtx); err != nil {
		return fmt.Errorf("shutdown %w",err)
	}
	log.Println("Server gracefully stopped")
	return nil
}