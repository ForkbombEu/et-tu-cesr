package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/ForkbombEu/et-tu-cesr/cesr"
	"github.com/spf13/cobra"
)

//go:embed schema/acdc/*/*
var schemaFiles embed.FS

func prettyPrint(events []cesr.Event, filterCreds bool) {
	for i, ev := range events {
		v, _ := ev.KED["v"].(string)
		if filterCreds && !strings.HasPrefix(v, "ACDC") {
			continue
		}
		fmt.Printf("\n### Event %d  (t=%v  sn=%v)\n", i+1, ev.KED["t"], ev.KED["s"])
		out, _ := json.MarshalIndent(ev.KED, "", "  ")
		fmt.Println(string(out))
		if ev.AttachBytes > 0 {
			fmt.Printf("• attachment bytes: %d\n", ev.AttachBytes)
		}
	}
}

func readCESRContent(path string, args []string) (string, error) {
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	// No path given - read from args or stdin
	if len(args) > 0 {
		return args[0], nil
	}

	// read from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("no input provided via -p, argument, or stdin")
	}
	return string(data), nil
}

func runDump(path string, args []string, credsOnly bool) error {
	content, err := readCESRContent(path, args)
	if err != nil {
		return err
	}

	events, err := cesr.ParseCESR(content)
	if err != nil {
		return err
	}
	prettyPrint(events, credsOnly)
	return nil
}

func runValidate(path string, args []string) error {
	content, err := readCESRContent(path, args)
	if err != nil {
		return err
	}
	events, err := cesr.ParseCESR(content)
	if err != nil {
		return err
	}
	subRoot, err := fs.Sub(schemaFiles, "schema/acdc")
	if err != nil {
		return err
	}
	v := cesr.NewValidator(subRoot)

	var errs []string
	valid := 0

	for idx, ev := range events {
		if err := v.ValidateCredential(ev.KED); err != nil {
			sn := ev.KED["s"]
			errs = append(errs, fmt.Sprintf("event %d (sn=%v) ⇒ %v", idx+1, sn, err))
			continue
		}
		if ver, _ := ev.KED["v"].(string); len(ver) >= 4 && ver[:4] == "ACDC" {
			valid++
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors:\n%s", strings.Join(errs, "\n"))
	}

	fmt.Printf("✅ %d credential bodies valid\n", valid)
	return nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "et-tu-cesr",
		Short: "Et tu, CESR? – stab CESR streams and reveal their JSON",
		Long: `Et tu, CESR? – stab CESR streams and reveal their JSON

Usage examples:
  et-tu-cesr dump -p file.cesr
  et-tu-cesr dump-credentials "cesr-content-string"
  echo "<cesr-content>" | et-tu-cesr validate-credentials
`,
	}

	var path string

	// dump command
	var dumpCmd = &cobra.Command{
		Use:   "dump [cesr-content]",
		Short: "pretty-print every event from CESR stream",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDump(path, args, false)
		},
	}
	dumpCmd.Flags().StringVarP(&path, "path", "p", "", "path to CESR file")

	// dump-credentials command
	var dumpCredsCmd = &cobra.Command{
		Use:   "dump-credentials [cesr-content]",
		Short: "pretty-print only ACDC credential bodies",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDump(path, args, true)
		},
	}
	dumpCredsCmd.Flags().StringVarP(&path, "path", "p", "", "path to CESR file")

	// validate-credentials command
	var validateCmd = &cobra.Command{
		Use:   "validate-credentials [cesr-content]",
		Short: "validate credential events in a CESR stream",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(path, args)
		},
	}
	validateCmd.Flags().StringVarP(&path, "path", "p", "", "path to CESR file")

	rootCmd.AddCommand(dumpCmd, dumpCredsCmd, validateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
