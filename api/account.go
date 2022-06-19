package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/martikan/simplebank/db/sqlc"
)

type listAccountsRequest struct { // TODO: Should be generic for all listing
	Page int32 `form:"page" binding:"min=1"`
	Size int32 `form:"size" binding:"min=5,max=50"`
	// q    string `form:"q"` // TODO: Should be implemented in the future for filtering
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=EUR USD"`
}

func (s *Server) listAccounts(ctx *gin.Context) {

	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListAccountsParams{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}

	accounts, err := s.store.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (s *Server) getAccount(ctx *gin.Context) {

	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (s *Server) createAccount(ctx *gin.Context) {

	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0, // TODO: It should be in service layer
	}

	account, err := s.store.CreateAccount(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
