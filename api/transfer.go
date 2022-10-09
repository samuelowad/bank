package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/samuelowad/bank/pkg/db/sqlc"
	"github.com/samuelowad/bank/pkg/token"
	"github.com/samuelowad/bank/pkg/util"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req transferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	fromAccount, valid := server.validateAccount(c, req.FromAccountID, req.Currency)

	if !valid {
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.PayLoad)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account not linked to current user")
		c.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}
	_, valid = server.validateAccount(c, req.ToAccountID, req.Currency)

	if !valid {
		return
	}

	arg := db.TransferTxParams{
		ToAccountID:   req.ToAccountID,
		FromAccountID: req.FromAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusCreated, result)

}

func (server *Server) validateAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	acc, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return acc, false

		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return acc, false
	}

	if acc.Currency != currency {
		err := fmt.Errorf("account[%d] currency not %s is not a valid currency", accountID, currency)
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return acc, false
	}

	return acc, true
}
