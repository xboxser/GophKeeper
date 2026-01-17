package model

type CredentialAPI struct {
	Login    string `json:"login" validate:"required"`
	Password []byte `json:"password" validate:"required"`
}

type Credential struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
