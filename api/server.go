package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/martikan/simplebank/db/sqlc"
	"github.com/martikan/simplebank/security"
	"github.com/martikan/simplebank/util"
)

// Server serves HTTP requests for the banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker security.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := security.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {

	router := gin.Default()

	// Open routes

	router.POST("/signin", s.signIn)

	router.POST("/users", s.createUser)

	// Authenticated routes

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.GET("/accounts", s.listAccounts)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.POST("/accounts", s.createAccount)

	authRoutes.POST("/transfers", s.createTransfer)

	s.router = router
}

// Start runs the HTTP server on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
