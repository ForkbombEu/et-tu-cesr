package cmd

import (
	"embed"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

// Execute runs the CLI.
func Execute(schemaFiles embed.FS) error {
	rootCmd = &cobra.Command{
		Use:   "et-tu-cesr",
		Short: "Et tu, CESR? – stab CESR streams and reveal their JSON",
		Long: `Et tu, CESR? – stab CESR streams and reveal their JSON

Usage examples:
  et-tu-cesr dump -p file.cesr
  et-tu-cesr dump-credentials "cesr-content-string"
  echo "<cesr-content>" | et-tu-cesr validate-credentials
`,
	}

	// Add subcommands
	rootCmd.AddCommand(newDumpCmd())
	rootCmd.AddCommand(newDumpCredsCmd())
	rootCmd.AddCommand(newValidateCmd(schemaFiles))

	return rootCmd.Execute()
}
