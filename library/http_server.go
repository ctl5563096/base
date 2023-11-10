package library

import (
	"context"
	"fmt"
	"net/http"
)

type HttpServer struct {
	*http.Server
}

func (s *HttpServer) Run() (err error) {
	go func() {
		if err = s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("http.Server.ListenAndServe: %w", err)
			return
		}
		err = nil
	}()
	return
}

func (s *HttpServer) Close() (err error)  {
	if err = s.Shutdown(context.TODO()); err != nil {
		err = fmt.Errorf("server shutdown: %w", err)
		return
	}
	return
}

func NewHttpServer(conf *HttpServerConfig,handler http.Handler)(server *HttpServer){
	server = &HttpServer{
		&http.Server{
			Addr: fmt.Sprintf(":%d", conf.Port),
		},
	}
	server.Handler = handler
	return server
}