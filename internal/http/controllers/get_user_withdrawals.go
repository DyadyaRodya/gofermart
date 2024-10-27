package controllers

import (
	"encoding/json"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	httpdto "github.com/DyadyaRodya/gofermart/internal/http/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors"
	"net/http"
)

type GetUserWithdrawalsController struct {
	interactor *interactors.GetWithdrawalsInteractor
}

func NewGetUserWithdrawalsController(interactor *interactors.GetWithdrawalsInteractor) *GetUserWithdrawalsController {
	return &GetUserWithdrawalsController{
		interactor: interactor,
	}
}

func (c *GetUserWithdrawalsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userInfo := ctx.Value(domainmodels.UserInfo{}).(*domainmodels.UserInfo)

	withdrawals, err := c.interactor.Handle(ctx, userInfo.UUID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	withdrawalsDTO := make([]*httpdto.WithdrawalDTO, 0, len(withdrawals))
	for _, withdrawal := range withdrawals {
		dto := httpdto.FromWithdrawalInfo(withdrawal)
		withdrawalsDTO = append(withdrawalsDTO, dto)
	}

	body, err := json.Marshal(withdrawalsDTO)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
