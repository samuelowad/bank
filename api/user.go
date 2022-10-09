package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/samuelowad/bank/pkg/db/sqlc"
	"github.com/samuelowad/bank/pkg/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3"`
	FullName string `json:"fullName" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, util.ErrorResponse(err))
				return
			}
		}
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashedPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Username:       req.Username,
		FullName:       req.FullName,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	res := newUserResponse(user)
	c.JSON(http.StatusCreated, res)

}

//user login
type loginUserRequest struct {
	Username string `json:"username" username:"required,username"`
	Password string `json:"password" binding:"required,min=3"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token" `
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	fmt.Println(req)

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	err = util.ComparePassword(req.Password, user.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	accessToken, err := server.TokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	c.JSON(http.StatusOK, rsp)

}
