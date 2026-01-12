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
