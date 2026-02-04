package handler

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Mock реализация command.CobraCommand
type mockCobraCommand struct{}

func (m *mockCobraCommand) GetCommands() []*cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}
	return []*cobra.Command{cmd}
}

func TestNewCli(t *testing.T) {
	version := "1.0.0"
	cli := NewCli(version)

	assert.NotNil(t, cli)
	assert.NotNil(t, cli.ClientCmd)
	assert.Equal(t, "todo", cli.ClientCmd.Use)
	assert.Equal(t, "Клиент для работы с gophkeeper", cli.ClientCmd.Short)
	assert.Equal(t, "Клиент для работы с gophkeeper", cli.ClientCmd.Long)
	assert.Equal(t, version, cli.ClientCmd.Version)
}

func TestAddCommand(t *testing.T) {
	cli := NewCli("")
	mockCmd := &mockCobraCommand{}

	initialCount := len(cli.ClientCmd.Commands())

	cli.AddCommand(mockCmd)
	afterCount := len(cli.ClientCmd.Commands())

	assert.Equal(t, initialCount+1, afterCount, "Команда должна быть добавлена")
	assert.Len(t, cli.ClientCmd.Commands(), 1)
	assert.Equal(t, "test", cli.ClientCmd.Commands()[0].Use)
}

func TestRun(t *testing.T) {
	cli := NewCli("")

	testCmd := &cobra.Command{
		Use: "testrun",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cli.ClientCmd.AddCommand(testCmd)
	cli.ClientCmd.SetArgs([]string{"testrun"})

	err := cli.Run()

	assert.NoError(t, err)
}
