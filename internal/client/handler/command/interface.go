package command

import "github.com/spf13/cobra"

type CobraCommand interface {
	GetCommands() []*cobra.Command
}
