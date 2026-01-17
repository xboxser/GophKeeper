package model

// Code - используется для передачи hash кода masterPass
type Code struct {
	Hash []byte `json:"hash" validate:"required"`
}
