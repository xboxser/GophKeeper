package handler

import (
	"gophkeeper/internal/client/handler/command"

	"github.com/spf13/cobra"
)

type Cli struct {
	ClientCmd *cobra.Command
}

func NewCli(buildVersion string) *Cli {

	clientCmd := &cobra.Command{
		Use:     "todo",
		Short:   "Клиент для работы с gophkeeper",
		Long:    "Клиент для работы с gophkeeper",
		Version: buildVersion,
	}

	return &Cli{ClientCmd: clientCmd}
}

// AddCommand - Собираем дерево команд
func (cli *Cli) AddCommand(cobraCommand command.CobraCommand) {
	commands := cobraCommand.GetCommands()
	cli.ClientCmd.AddCommand(commands...)
}

// Run - запускает cli клиент
func (cli *Cli) Run() error {
	return cli.ClientCmd.Execute()
}
