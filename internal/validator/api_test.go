package validator

import (
	"gophkeeper/internal/model"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateModelUserAPI(t *testing.T) {
	if Validate == nil {
		Init()
	}

	tests := []struct {
		name     string
		user     model.APIUser
		valError bool
	}{
		{name: "valid parameters", user: model.APIUser{Login: "qwerty", Password: "pass"}, valError: false},
		{name: "empty parameters", user: model.APIUser{Login: "", Password: ""}, valError: true},
		{name: "nil password", user: model.APIUser{Login: "qwerty"}, valError: true},
		{name: "nil login", user: model.APIUser{Password: "pass"}, valError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModelUserAPI(tt.user)
			if !tt.valError {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestValidateModelRegistrationUserAPI(t *testing.T) {
	if Validate == nil {
		Init()
	}

	tests := []struct {
		name     string
		user     model.APIUser
		valError bool
	}{
		{name: "valid parameters", user: model.APIUser{Login: "qwerty", Password: "pass", Code: []byte{1}}, valError: false},
		{name: "empty parameters", user: model.APIUser{Login: "", Password: "", Code: []byte{}}, valError: true},
		{name: "empty code", user: model.APIUser{Login: "qwe", Password: "qwe", Code: []byte{}}, valError: true},
		{name: "nil code", user: model.APIUser{Login: "qwe", Password: "qwe", Code: nil}, valError: true},
		{name: "nil password", user: model.APIUser{Login: "qwerty"}, valError: true},
		{name: "nil login", user: model.APIUser{Password: "pass"}, valError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModelRegistrationUserAPI(tt.user)
			if !tt.valError {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestValidateModelCredentialAPI(t *testing.T) {
	if Validate == nil {
		Init()
	}

	tests := []struct {
		name       string
		credential model.Credential
		valError   bool
	}{
		{name: "valid parameters", credential: model.Credential{Login: "qwerty", Password: "pass"}, valError: false},
		{name: "empty parameters", credential: model.Credential{Login: "", Password: ""}, valError: true},
		{name: "nil password", credential: model.Credential{Login: "qwerty"}, valError: true},
		{name: "nil login", credential: model.Credential{Password: "pass"}, valError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModelCredentialAPI(tt.credential)
			if !tt.valError {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}

func TestValidateModelCardAPI(t *testing.T) {
	if Validate == nil {
		Init()
	}

	tests := []struct {
		name     string
		card     model.Card
		valError bool
	}{
		{name: "valid parameters", card: model.Card{Title: "title", Number: "number", Expiry: "Expiry", CVV: "CVV", CardHolder: "CardHolder"}, valError: false},
		{name: "empty parameters", card: model.Card{Title: "", Number: "", Expiry: "", CVV: "", CardHolder: "CardHolder"}, valError: true},
		{name: "nil Title", card: model.Card{Title: "", Number: "number", Expiry: "Expiry", CVV: "CVV", CardHolder: "CardHolder"}, valError: true},
		{name: "nil Number", card: model.Card{Title: "title", Number: "", Expiry: "Expiry", CVV: "CVV", CardHolder: "CardHolder"}, valError: true},
		{name: "nil Expiry", card: model.Card{Title: "title", Number: "number", Expiry: "", CVV: "CVV", CardHolder: "CardHolder"}, valError: true},
		{name: "nil CVV", card: model.Card{Title: "title", Number: "number", Expiry: "Expiry", CVV: "", CardHolder: "CardHolder"}, valError: true},
		{name: "nil CardHolder", card: model.Card{Title: "title", Number: "number", Expiry: "Expiry", CVV: "CVV", CardHolder: ""}, valError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModelCardAPI(tt.card)
			if !tt.valError {
				require.NoError(t, err)
				return
			}
			assert.Error(t, err)
		})
	}
}
