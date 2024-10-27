package controllers

import (
	"encoding/json"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	httpdto "github.com/DyadyaRodya/gofermart/internal/http/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors"
	"net/http"
)

type GetBalanceController struct {
	interactor *interactors.GetBalanceInteractor
}

func NewGetBalanceController(interactor *interactors.GetBalanceInteractor) *GetBalanceController {
	return &GetBalanceController{
		interactor: interactor,
	}
}

func (c *GetBalanceController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	userInfo := ctx.Value(domainmodels.UserInfo{}).(*domainmodels.UserInfo)

	balanceInfo, err := c.interactor.Handle(ctx, userInfo.UUID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	dto := httpdto.FromBalanceInfo(balanceInfo)
	body, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
