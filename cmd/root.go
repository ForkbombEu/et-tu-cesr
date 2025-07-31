package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"
var rootCmd = &cobra.Command{
	Use:   "et-tu-cesr",
	Short: "Et tu, CESR? – stab CESR streams and reveal their JSON",
	Long: `Et tu, CESR? – stab CESR streams and reveal their JSON

Usage examples:
  et-tu-cesr dump -p file.cesr
  et-tu-cesr dump-credentials "cesr-content-string"
  echo "<cesr-content>" | et-tu-cesr validate-credentials
`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip printing if this is the version command or version flag
		if cmd.Name() == "version" {
			return
		}
		// Check if --version flag is set (in case of root command)
		versionFlag := cmd.Flags().Lookup("version")
		if versionFlag != nil && versionFlag.Changed {
			return
		}

		// Print version info on stderr before running any command
		fmt.Fprintln(os.Stderr, "et-tu-cesr", version)
	},
}

func Execute(schemaFiles embed.FS) error {
	// Add subcommands
	rootCmd.AddCommand(newDumpCmd())
	rootCmd.AddCommand(newDumpCredsCmd())
	rootCmd.AddCommand(newValidateCmd(schemaFiles))
	rootCmd.AddCommand(newValidateParsedCmd(schemaFiles))
	return rootCmd.Execute()
}
