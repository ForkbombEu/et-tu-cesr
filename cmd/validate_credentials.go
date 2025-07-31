package cmd

import (
	"embed"

	"github.com/ForkbombEu/et-tu-cesr/internal"

	"github.com/spf13/cobra"
)

func newValidateCmd(schemaFiles embed.FS) *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "validate-credentials [cesr-content]",
		Short: "validate credential events in a CESR stream",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.RunValidate(path, args, schemaFiles)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", "", "path to CESR file")
	return cmd
}
