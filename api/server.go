package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/martikan/simplebank/db/sqlc"
)

// Server serves HTTP requests for the banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {

	server := &Server{store: store}

	router := gin.Default()

	router.GET("/accounts", server.listAccounts)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
