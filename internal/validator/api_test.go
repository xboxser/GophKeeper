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
	password := []byte("pass")

	tests := []struct {
		name       string
		credential model.CredentialAPI
		valError   bool
	}{
		{name: "valid parameters", credential: model.CredentialAPI{Login: "qwerty", Password: password}, valError: false},
		{name: "empty parameters", credential: model.CredentialAPI{Login: "", Password: []byte{}}, valError: true},
		{name: "nil password", credential: model.CredentialAPI{Login: "qwerty"}, valError: true},
		{name: "nil login", credential: model.CredentialAPI{Password: password}, valError: true},
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

	number := []byte("number")
	expiry := []byte("expiry")
	cvv := []byte("cvv")
	cartHolder := []byte("cartHolder")

	tests := []struct {
		name     string
		card     model.CardAPI
		valError bool
	}{
		{name: "valid parameters", card: model.CardAPI{Title: "title", Number: number, Expiry: expiry, CVV: cvv, CardHolder: cartHolder}, valError: false},
		{name: "empty parameters", card: model.CardAPI{Title: "", Number: []byte{}, Expiry: []byte{}, CVV: []byte{}, CardHolder: []byte{}}, valError: true},
		{name: "nil Title", card: model.CardAPI{Title: "", Number: number, Expiry: expiry, CVV: cvv, CardHolder: cartHolder}, valError: true},
		{name: "nil Number", card: model.CardAPI{Title: "title", Number: []byte{}, Expiry: expiry, CVV: cvv, CardHolder: cartHolder}, valError: true},
		{name: "nil Expiry", card: model.CardAPI{Title: "title", Number: number, Expiry: []byte{}, CVV: cvv, CardHolder: cartHolder}, valError: true},
		{name: "nil CVV", card: model.CardAPI{Title: "title", Number: number, Expiry: expiry, CVV: []byte{}, CardHolder: cartHolder}, valError: true},
		{name: "nil CardHolder", card: model.CardAPI{Title: "title", Number: number, Expiry: expiry, CVV: cvv, CardHolder: []byte{}}, valError: true},
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
