package model

import "errors"

type TypeCounter string

const (
	// CounterCard - тип счетчика для карт
	CounterCard TypeCounter = "card"
	// CounterCredential - тип счетчика для учетных данных
	CounterCredential TypeCounter = "cred"
)

var (
	ErrCounterAddCard       = errors.New("error create counter for card")
	ErrCounterAddCredential = errors.New("error create counter for credential")
)
