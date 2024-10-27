package controllers

import (
	"encoding/json"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	httpdto "github.com/DyadyaRodya/gofermart/internal/http/dto"
	"github.com/DyadyaRodya/gofermart/internal/interactors"
	"net/http"
)

type GetUserOrdersController struct {
	interactor *interactors.GetOrdersInteractor
}

func NewGetUserOrdersController(interactor *interactors.GetOrdersInteractor) *GetUserOrdersController {
	return &GetUserOrdersController{
		interactor: interactor,
	}
}

func (c *GetUserOrdersController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userInfo := ctx.Value(domainmodels.UserInfo{}).(*domainmodels.UserInfo)

	orders, err := c.interactor.Handle(ctx, userInfo.UUID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ordersDTO := make([]*httpdto.OrderDTO, 0, len(orders))
	for _, order := range orders {
		dto := httpdto.FromOrderInfo(order)
		ordersDTO = append(ordersDTO, dto)
	}

	body, err := json.Marshal(ordersDTO)
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
