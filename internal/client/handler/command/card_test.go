package command

import (
	"gophkeeper/mocks/client/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardHandler_GetCommands(t *testing.T) {
	mockCardService := &service.MockCardService{}
	cardHandler := NewCardHandler(mockCardService, nil)

	commands := cardHandler.GetCommands()

	assert.Len(t, commands, 4)

	commandUses := make([]string, 0, 4)
	for _, cmd := range commands {
		commandUses = append(commandUses, cmd.Use)
	}

	assert.Contains(t, commandUses, "card-add")
	assert.Contains(t, commandUses, "card-get")
	assert.Contains(t, commandUses, "card-delete")
	assert.Contains(t, commandUses, "card-update")
}
