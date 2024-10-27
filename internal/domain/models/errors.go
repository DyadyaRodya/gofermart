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

var ErrLoginValidation = errors.Join(ErrInteractor, errors.New("login validation error"))
var ErrLoginTooLong = errors.Join(ErrLoginValidation, errors.New("login too long"))
var ErrLoginTooShort = errors.Join(ErrLoginValidation, errors.New("login too short"))
var ErrLoginChars = errors.Join(ErrLoginValidation, errors.New("login contains incorrect chars"))
var ErrLoginTaken = errors.Join(ErrInteractor, errors.New("login already taken error"))
var ErrWrongCredentials = errors.Join(ErrInteractor, errors.New("wrong credentials"))
var ErrPasswordComplexity = errors.Join(ErrInteractor, errors.New("password complexity error"))
var ErrUserNotFound = errors.Join(ErrInteractor, errors.New("user not found"))

// Order errors

var ErrOrderExists = errors.Join(ErrInteractor, errors.New("order already exists"))
var ErrSameUserOrderExists = errors.Join(ErrInteractor, errors.New("order already exists by same user"))
var ErrOrderNumberInvalid = errors.Join(ErrInteractor, errors.New("order number invalid"))

// Withdraw errors

var ErrNotEnoughPointsToWithdraw = errors.Join(ErrInteractor, errors.New("not enough points to withdraw"))
var ErrWithdrawExists = errors.Join(ErrInteractor, errors.New("withdraw exists"))

// Accrual

var ErrAccrualOrderNotRegistered = errors.Join(ErrInteractor, errors.New("order not registered in accrual system"))
