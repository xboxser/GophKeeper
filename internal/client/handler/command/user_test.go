package command

import (
	service "gophkeeper/mocks/client/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetCommands(t *testing.T) {
	mockUserService := &service.MockUserService{}
	userHandler := NewUserHandler(mockUserService)

	commands := userHandler.GetCommands()

	assert.Len(t, commands, 3)

	commandUses := make([]string, 0, 3)
	for _, cmd := range commands {
		commandUses = append(commandUses, cmd.Use)
	}
	assert.Contains(t, commandUses, "registration")
	assert.Contains(t, commandUses, "login")
	assert.Contains(t, commandUses, "master")
}
