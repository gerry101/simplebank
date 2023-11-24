package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/util"
)

type Server struct {
	store db.Store
	router *gin.Engine
	config util.Config
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store: store,
		config: config,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err = v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Fatal("could not register custom currency validator: ", err)
			return nil, err
		}
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authGroup := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authGroup.GET("/accounts", server.listAccounts)
	authGroup.GET("/accounts/:id", server.getAccount)
	authGroup.POST("/accounts", server.createAccount)
	authGroup.DELETE("/accounts/:id", server.deleteAccount)

	authGroup.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
