package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ForkbombEu/et-tu-cesr/cesr"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Et tu, CESR? – stab CESR streams and reveal their JSON

Usage:
    et-tu-cesr dump <file.cesr>              pretty‑print every event
    et-tu-cesr dump-credentials <file.cesr>  pretty‑print only ACDC credential bodies
    et-tu-cesr help                          show this message
`)
}

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

func runDump(file string, credsOnly bool) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	events, err := cesr.ParseCESR(string(data))
	if err != nil {
		return err
	}
	prettyPrint(events, credsOnly)
	return nil
}

// runValidate validates every credential event in a .cesr file
// using the super‑simple Validator (no caching, schema chosen by filename).
// It prints a summary and returns an error on first failure.
func runValidate(file string) error {
	// 1 – read the CESR stream
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	events, err := cesr.ParseCESR(string(data))
	if err != nil {
		return err
	}

	// 2 – create a lightweight validator that points at your schema folder
	v := cesr.NewValidator("schema/acdc") // adjust if your path differs

	// 3 – walk events
	valid := 0
	for idx, ev := range events {
		if err := v.ValidateCredential(ev.KED); err != nil {
			sn := ev.KED["s"]
			return fmt.Errorf("%s: event %d (sn=%v) ⇒ %v", file, idx+1, sn, err)
		}
		// count only credential bodies (v starts with ACDC)
		if ver, _ := ev.KED["v"].(string); len(ver) >= 4 && ver[:4] == "ACDC" {
			valid++
		}
	}

	fmt.Printf("✅ %d credential bodies valid in %s\n", valid, file)
	return nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}

	switch cmd := flag.Arg(0); cmd {
	case "dump":
		if flag.NArg() != 2 {
			usage()
			os.Exit(1)
		}
		if err := runDump(flag.Arg(1), false); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

	case "dump-credentials":
		if flag.NArg() != 2 {
			usage()
			os.Exit(1)
		}
		if err := runDump(flag.Arg(1), true); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

	case "validate-credentials":
		if flag.NArg() != 2 {
			usage()
			os.Exit(1)
		}
		if err := runValidate(flag.Arg(1)); err != nil {
			fmt.Fprintln(os.Stderr, "validation failed:", err)
			os.Exit(1)
		}

	case "help":
		usage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", cmd)
		usage()
		os.Exit(1)
	}
}
