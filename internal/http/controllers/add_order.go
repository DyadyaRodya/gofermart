package controllers

import (
	"errors"
	domainmodels "github.com/DyadyaRodya/gofermart/internal/domain/models"
	"github.com/DyadyaRodya/gofermart/internal/interactors"
	"io"
	"net/http"
)

type AddOrderController struct {
	interactor *interactors.AddOrderInteractor
}

func NewAddOrderController(interactor *interactors.AddOrderInteractor) *AddOrderController {
	return &AddOrderController{
		interactor: interactor,
	}
}

func (c *AddOrderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "provide order number in text/plain body", http.StatusBadRequest)
		return
	}

	orderNumber := string(body)

	userInfo := ctx.Value(domainmodels.UserInfo{}).(*domainmodels.UserInfo)

	err = c.interactor.Handle(ctx, userInfo.UUID, orderNumber)
	if err != nil {
		switch {
		case errors.Is(err, domainmodels.ErrSameUserOrderExists):
			w.WriteHeader(http.StatusOK)
		case errors.Is(err, domainmodels.ErrOrderExists):
			w.WriteHeader(http.StatusConflict)
		case errors.Is(err, domainmodels.ErrOrderNumberInvalid):
			w.WriteHeader(http.StatusUnprocessableEntity)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
