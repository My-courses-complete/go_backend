package api

import (
	"github.com/My-courses-complete/go_backend.git/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	store *pgxpool.Pool
	router *gin.Engine
	*db.Queries
}

func NewServer(store *pgxpool.Pool) *Server {
	server := &Server{
		store: store,
	}
	server.Queries = db.New()
	
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (server *Server) Run(addr string) error {
    return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{
        "error": err.Error(),
    }
}