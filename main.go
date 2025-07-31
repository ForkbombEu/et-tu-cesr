package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/ForkbombEu/et-tu-cesr/cmd"
)

//go:embed schema/acdc/*/*
var schemaFiles embed.FS

func main() {
	if err := cmd.Execute(schemaFiles); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
