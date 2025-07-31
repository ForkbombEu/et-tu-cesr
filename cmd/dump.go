package cmd

import (
	"github.com/ForkbombEu/et-tu-cesr/internal"

	"github.com/spf13/cobra"
)

func newDumpCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "dump [cesr-content]",
		Short: "pretty-print every event from CESR stream",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.RunDump(path, args, false)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "path to CESR file")
	return cmd
}
