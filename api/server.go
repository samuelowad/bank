package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/samuelowad/bank/api/middleware"
	db "github.com/samuelowad/bank/pkg/db/sqlc"
	"github.com/samuelowad/bank/pkg/token"
	"github.com/samuelowad/bank/pkg/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	TokenMaker token.Maker
	Router     *gin.Engine
}

//NewServer creates a new HTTP server and routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenMaker :%w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		TokenMaker: tokenMaker}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", middleware.ValidCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(server.TokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	//transaction
	authRoutes.POST("/transfer", server.createTransfer)

	//user
	router.POST("/users", server.createUser)

	router.POST("/users/login", server.loginUser)

	//	add router to router table
	server.Router = router
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
