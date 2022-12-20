package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/server/handler"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	//mux      http.Handler

	// Handlers
	ah *handler.AuthHandler
	//ms       mid.MiddlewareService

	//mux      *pat.PatternServeMux
	//mux      *http.ServeMux
}

func NewServer(ctx context.Context, host string, port uint, authHanlder *handler.AuthHandler) Server { //, hr *handlers.HandlerRepo,m mid.MiddlewareService) Server { //(context.Context, Server) {
	//func New(ctx context.Context, host string, port uint, router http.Handler, hr *handlers.HandlerRepo) Server { //(context.Context, Server) {
	srv := Server{
		httpAddr: fmt.Sprintf(host + ":" + fmt.Sprint(port)),
		engine:   gin.New(),
		//mux:      http.NewServeMux(),
		//mux: pat.New(),

		// Handlers
		ah: authHanlder,
		//ms: m,
	}

	//return serverContext(ctx), srv
	return srv

}

func (s *Server) Run(ctx context.Context) error {
	log.Printf("Listening on %s\n", s.httpAddr)
	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	s.registerRoutes()

	return srv.ListenAndServe()
	//return http.ListenAndServe(s.httpAddr, s.engine)
}

func (s *Server) registerRoutes() {
	fmt.Println("Engine routes ...")
	s.engine.GET("/health", s.ah.CheckHandler())
	s.engine.GET("/users", s.ah.FindUsers())
	s.engine.GET("/usersc", s.ah.FindUsersC)

}
