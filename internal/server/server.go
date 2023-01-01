package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	auth_handler "github.com/krlspj/go-jwt-auth/internal/auth/platform/server/handler"
	authz_handler "github.com/krlspj/go-jwt-auth/internal/authz/platform/server/handler"
)

const (
	SERVER_GIN = "gin"
	SERVER_GO  = "native"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine

	//mux      http.Handler

	// Handlers
	ah *auth_handler.AuthHandler
	az *authz_handler.AuthzMidd
	//ms       mid.MiddlewareService

	//mux      *pat.PatternServeMux
	//mux      *http.ServeMux
}

func NewServer(
	ctx context.Context,
	host string,
	port uint,
	authHanlder *auth_handler.AuthHandler,
	authzMidd *authz_handler.AuthzMidd,

) Server { //, hr *handlers.HandlerRepo,m mid.MiddlewareService) Server { //(context.Context, Server) {
	//func New(ctx context.Context, host string, port uint, router http.Handler, hr *handlers.HandlerRepo) Server { //(context.Context, Server) {
	srv := Server{
		httpAddr: fmt.Sprintf(host + ":" + fmt.Sprint(port)),
		engine:   gin.Default(),
		//engine: gin.New(),
		//mux:      http.NewServeMux(),
		//mux: pat.New(),

		// Handlers
		ah: authHanlder,
		az: authzMidd,
		//ms: m,
	}

	//return serverContext(ctx), srv
	return srv

}

func (s *Server) Run(ctx context.Context, serverType string) error {
	log.Printf("Listening on %s\n", s.httpAddr)

	s.registerRoutes()

	switch serverType {
	case SERVER_GIN:
		port := "60002"
		log.Println("\033[00;32m[NOTICE] Starting Gin server on port:", port, "\033[m")
		return s.engine.Run(":60002")

	case SERVER_GO:
		log.Println("\033[00;32m[NOTICE] Native server selected.\033[m")
		//		srv := &http.Server{
		//			Addr:    s.httpAddr,
		//			Handler: s.engine,
		//		}
		//		log.Println("\033[00;32m[NOTICE] Starting native server:", s.httpAddr, "\033[m")
		//		return srv.ListenAndServe()
		fallthrough

	default:
		srv := &http.Server{
			Addr:    s.httpAddr,
			Handler: s.engine,
		}
		log.Println("\033[00;32m[NOTICE] Starting native http.Server:", s.httpAddr, "\033[m")
		return srv.ListenAndServe()
	}

	//return http.ListenAndServe(s.httpAddr, s.engine)
}

func (s *Server) registerRoutes() {
	fmt.Println("Engine routes ...")
	s.engine.GET("/health", s.ah.CheckHandler())
	s.engine.POST("/login", s.ah.Login)
	s.engine.POST("/register", s.ah.CreateNewUser)

	usersGroup := s.engine.Group("/users")
	usersGroup.Use(s.az.VeriyToken)

	usersGroup.GET("/", s.az.OnlyAdmin(), s.ah.FindUsers())
	usersGroup.GET("/c", s.ah.FindUsersC)
	usersGroup.GET("/:id", s.ah.FindUser)
	usersGroup.GET("/by_field", s.ah.FindUserByField)
	usersGroup.PUT("/:id", s.ah.PutUser)
	usersGroup.PATCH("/:id", s.ah.PatchUser)

	configGroup := s.engine.Group("/v1/config")
	configGroup.GET("/:id")
	configGroup.PUT("/:id")
	configGroup.PATCH("/:id")

}
