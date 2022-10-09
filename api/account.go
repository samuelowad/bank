package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/samuelowad/bank/pkg/db/sqlc"
	"github.com/samuelowad/bank/pkg/token"
	"github.com/samuelowad/bank/pkg/util"
	"net/http"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(c *gin.Context) {
	var req createAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.PayLoad)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(c, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var req getAccountRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := server.store.GetAccount(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return

		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.PayLoad)
	if account.Owner != authPayload.Username {
		err := errors.New("account not belonging to authenticated user")
		c.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(c *gin.Context) {
	var req listAccountRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.PayLoad)

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := server.store.ListAccounts(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return

		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)

}
