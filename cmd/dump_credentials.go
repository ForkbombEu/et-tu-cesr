package cmd

import (
	"github.com/ForkbombEu/et-tu-cesr/internal"

	"github.com/spf13/cobra"
)

func newDumpCredsCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "dump-credentials [cesr-content]",
		Short: "pretty-print only ACDC credential bodies",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.RunDump(path, args, true)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "path to CESR file")
	return cmd
}
