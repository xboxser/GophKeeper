package command

import (
	service "gophkeeper/mocks/client/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileHandler_GetCommands(t *testing.T) {
	mockFileService := &service.MockFileService{}
	fileHandler := NewFileHandler(mockFileService, nil)

	commands := fileHandler.GetCommands()
	assert.Len(t, commands, 4)

	commandUses := make([]string, 0, 4)
	for _, cmd := range commands {
		commandUses = append(commandUses, cmd.Use)
	}

	assert.Contains(t, commandUses, "file-add")
	assert.Contains(t, commandUses, "file-list")
	assert.Contains(t, commandUses, "file-get")
	assert.Contains(t, commandUses, "file-delete")
}
