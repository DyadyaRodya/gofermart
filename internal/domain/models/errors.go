package models

import (
	"errors"
)

// internal unexpected errors

var ErrInternalServer = errors.New("internal server error")
var ErrPasswordHashGeneration = errors.Join(ErrInternalServer, errors.New("password hash generation error"))
var ErrSaltGeneration = errors.Join(ErrInternalServer, errors.New("salt generation error"))

// interactor errors

var ErrInteractor = errors.New("interactor error")

// User errors

var ErrLoginTaken = errors.Join(ErrInteractor, errors.New("login already taken error"))
var ErrWrongCredentials = errors.Join(ErrInteractor, errors.New("wrong credentials"))
var ErrPasswordComplexity = errors.Join(ErrInteractor, errors.New("password complexity error"))
var ErrUserNotFound = errors.Join(ErrInteractor, errors.New("user not found"))

// Order errors

var ErrOrderExists = errors.Join(ErrInteractor, errors.New("order already exists"))
var ErrOrderNumberInvalid = errors.Join(ErrInteractor, errors.New("order number invalid"))

// Withdraw errors

var ErrNotEnoughPointsToWithdraw = errors.Join(ErrInteractor, errors.New("not enough points to withdraw"))
var ErrWithdrawExists = errors.Join(ErrInteractor, errors.New("withdraw exists"))

// Accrual

var ErrAccrualOrderNotRegistered = errors.Join(ErrInteractor, errors.New("order not registered in accrual system"))
