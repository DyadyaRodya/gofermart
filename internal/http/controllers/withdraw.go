package controllers

import (
	"encoding/json"
	"errors"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	httpdto "github.com/DyadyaRodya/gofermart/internal/http/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors"
	"net/http"
)

type WithdrawController struct {
	interactor *interactors.WithdrawInteractor
}

func NewWithdrawController(interactor *interactors.WithdrawInteractor) *WithdrawController {
	return &WithdrawController{
		interactor: interactor,
	}
}

func (c *WithdrawController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var withdrawDTO httpdto.WithdrawDTO
	err := json.NewDecoder(r.Body).Decode(&withdrawDTO)
	if err != nil {
		http.Error(w, "You need to provide order and sum in JSON body", http.StatusBadRequest)
		return
	}
	orderNumber := withdrawDTO.OrderNumber
	sum := withdrawDTO.ConvertSum()

	userInfo := ctx.Value(domainmodels.UserInfo{}).(*domainmodels.UserInfo)

	err = c.interactor.Handle(ctx, userInfo.UUID, orderNumber, sum)
	if err != nil {
		switch {
		case errors.Is(err, domainmodels.ErrOrderNumberInvalid):
			w.WriteHeader(http.StatusUnprocessableEntity)
		case errors.Is(err, domainmodels.ErrNotEnoughPointsToWithdraw):
			w.WriteHeader(http.StatusPaymentRequired)
		case errors.Is(err, domainmodels.ErrWithdrawExists):
			w.WriteHeader(http.StatusOK) // idempotent
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
