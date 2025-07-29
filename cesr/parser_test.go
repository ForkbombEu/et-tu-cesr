// parser_test.go
package cesr

import (
	"os"
	"path/filepath"
	"testing"
)

// helper: assert no error
func must(t *testing.T, err error) {
	if err != nil {
		t.Helper()
		t.Fatalf("unexpected error: %v", err)
	}
}

// Test the bodyLength helper with a known header.
func TestBodyLength(t *testing.T) {
	got, err := bodyLength("KERI10JSON000249_")
	must(t, err)
	if got != 0x0249 {
		t.Fatalf("want 0x0249, got %d", got)
	}
}

// Walk every *.cesr under ./samples and parse it.
func TestSampleStreams(t *testing.T) {
	glob := filepath.Join("..", "samples", "*.cesr") // adjust if path differs
	files, err := filepath.Glob(glob)
	must(t, err)
	if len(files) == 0 {
		t.Fatalf("no sample files found at %s", glob)
	}

	for _, f := range files {
		f := f // capture
		t.Run(filepath.Base(f), func(t *testing.T) {
			data, err := os.ReadFile(f)
			must(t, err)

			events, err := ParseCESR(string(data))
			must(t, err)
			if len(events) == 0 {
				t.Fatalf("parsed zero events")
			}

			// quick sanity checks on the first event
			e0 := events[0].KED
			if e0["v"] == nil || e0["t"] == nil {
				t.Fatalf("missing mandatory fields in first event: %#v", e0)
			}
		})
	}
}
