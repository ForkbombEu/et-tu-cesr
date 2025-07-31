package internal

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
)

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

func readContent(path string, args []string) (string, error) {
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	if len(args) > 0 {
		return args[0], nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", errors.New("no input provided via -p, argument, or stdin")
	}
	return string(data), nil
}

func RunDump(path string, args []string, credsOnly bool) error {
	content, err := readContent(path, args)
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

func RunValidate(path string, args []string, schemaFiles embed.FS) error {
	content, err := readContent(path, args)
	if err != nil {
		return err
	}
	events, err := cesr.ParseCESR(content)
	if err != nil {
		return err
	}
	return validateEvents(events, schemaFiles)
}
func RunValidateParsedJSON(path string, args []string, schemaFiles embed.FS) error {
	jsonStr, err := readContent(path, args)
	if err != nil {
		return err
	}

	var events []cesr.Event
	if err := json.Unmarshal([]byte(jsonStr), &events); err != nil {
		return fmt.Errorf("failed to unmarshal events JSON: %w", err)
	}

	return validateEvents(events, schemaFiles)
}

func validateEvents(events []cesr.Event, schemaFiles embed.FS) error {
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
