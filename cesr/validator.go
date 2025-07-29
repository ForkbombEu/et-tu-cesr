package cesr

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validator walks a short list of schema folders – newest first – and
// validates a single ACDC credential event.  No caching involved.
type Validator struct{ dirs []string }

// NewValidator builds an *absolute* search list:
//
//	<root>/2023/acdc/   (current spec)
//	<root>/2022/acdc/   (legacy)
//
// Pass the root once (e.g. "schema") and forget about cwd quirks.
func NewValidator(root string) *Validator {
	abs, err := filepath.Abs(root)
	if err != nil {
		// fallback to the original string – Compile will fail clearly later
		abs = root
	}
	return &Validator{
		dirs: []string{
			filepath.Join(abs, "2023"),
			filepath.Join(abs, "2022"),
		},
	}
}

// -----------------------------------------------------------------------------
// tiny heuristic → schema filename
// -----------------------------------------------------------------------------
func chooseFile(ked map[string]interface{}) (string, error) {
	a, ok := ked["a"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("missing \"a\" section")
	}
	switch {
	case a["LEI"] != nil && a["f"] != nil:
		return "legal-entity-engagement-context-role-vLEI-credential.json", nil
	case a["officialRole"] != nil:
		return "legal-entity-official-organizational-role-vLEI-credential.json", nil
	case a["LEI"] != nil:
		return "legal-entity-vLEI-credential.json", nil
	default:
		return "", fmt.Errorf("unrecognised credential type")
	}
}

// -----------------------------------------------------------------------------
// ValidateCredential – try every version folder until one compiles *and*
// accepts the credential.  Return a composite error otherwise.
// -----------------------------------------------------------------------------
func (v *Validator) ValidateCredential(ked map[string]interface{}) error {
	// ignore non‑ACDC events
	if ver, _ := ked["v"].(string); !strings.HasPrefix(ver, "ACDC") {
		return nil
	}

	fname, err := chooseFile(ked)
	if err != nil {
		return err
	}

	var attempts []string
	for _, dir := range v.dirs {
		path := filepath.Join(dir, fname)
		uri := "file://" + filepath.ToSlash(path) // absolute file URI

		c := jsonschema.NewCompiler()
		sch, err := c.Compile(uri)
		if err != nil {
			attempts = append(attempts,
				fmt.Sprintf("%s (compile error: %v)", uri, err))
			continue
		}
		if err = sch.Validate(ked); err == nil {
			return nil // 🎉 validated successfully
		}
		attempts = append(attempts,
			fmt.Sprintf("%s (validation error: %v)", uri, err))
	}

	return fmt.Errorf("no schema accepted credential; tried:\n  %s",
		strings.Join(attempts, "\n  "))
}
