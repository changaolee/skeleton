package util

import "github.com/spf13/cobra"

type CommandGroup struct {
	Message  string
	Commands []*cobra.Command
}

type CommandGroups []CommandGroup

func (g CommandGroups) AddTo(c *cobra.Command) {
	for _, group := range g {
		c.AddCommand(group.Commands...)
	}
}

func (g CommandGroups) Has(c *cobra.Command) bool {
	for _, group := range g {
		for _, command := range group.Commands {
			if command == c {
				return true
			}
		}
	}
	return false
}
