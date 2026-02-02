package command

import (
	service "gophkeeper/mocks/client/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialHandler_GetCommands(t *testing.T) {
	mockCredentialService := &service.MockCredentialService{}
	credentialHandler := NewCredentialHandler(mockCredentialService, nil)

	commands := credentialHandler.GetCommands()

	assert.Len(t, commands, 4)

	commandUses := make([]string, 0, 4)
	for _, cmd := range commands {
		commandUses = append(commandUses, cmd.Use)
	}

	assert.Contains(t, commandUses, "credential-add")
	assert.Contains(t, commandUses, "credential-get")
	assert.Contains(t, commandUses, "credential-del")
	assert.Contains(t, commandUses, "credential-update")
}
