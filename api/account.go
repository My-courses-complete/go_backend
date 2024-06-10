package api

import (
	"database/sql"
	"net/http"

	"github.com/My-courses-complete/go_backend.git/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}


func (s *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err!= nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
        return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
        Currency: req.Currency,
		Balance: 0,
	}

	account, err := s.Queries.CreateAccount(c.Request.Context(), s.store, arg)
	if err!= nil {
        c.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

	c.JSON(http.StatusCreated, account)
}

type getAccountParams struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) getAccount(c *gin.Context) {
	var req getAccountParams
    if err := c.ShouldBindUri(&req); err!= nil {
        c.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    account, err := s.Queries.GetAccount(c.Request.Context(), s.store, req.ID)
    if err!= nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
        c.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    c.JSON(http.StatusOK, account)
}

type listAccountsParams struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccounts(c *gin.Context) {
	var req listAccountsParams
	if err := c.ShouldBindQuery(&req); err!= nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := s.Queries.ListAccounts(c.Request.Context(), s.store, arg)
	if err!= nil {
        c.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

	c.JSON(http.StatusOK, accounts)
}