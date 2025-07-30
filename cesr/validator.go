package cesr

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validator walks a short list of schema folders – newest first – and
// validates a single ACDC credential event. No caching involved.
// Validator walks through schema directories (newest first) using fs.FS
type Validator struct {
	root     fs.FS             // Root directory for schema files
	versions []versionCompiler // Each version has a name and a compiler
}

type versionCompiler struct {
	name     string
	compiler *jsonschema.Compiler
}
type MissingSchemaError struct{}

func (MissingSchemaError) Error() string {
	return "the schema file is missing or unreadable"
}

// NewValidator creates a Validator with embedded schema directories
func NewValidator(root fs.FS) *Validator {
	versions := []versionCompiler{}

	for _, year := range []string{"2023", "2022"} {
		_, err := fs.Stat(root, year)
		if err != nil {
			continue // Skip missing version
		}
		versions = append(versions, versionCompiler{
			name: year,
		})
	}

	return &Validator{
		root:     root,
		versions: versions,
	}
}

// createCompiler loads all JSON schemas in an FS into a compiler
func (vc *versionCompiler) updateCompiler(fsys fs.FS, version string, path string) error {
	if vc.compiler == nil {
		vc.compiler = jsonschema.NewCompiler()
	}
	// Read schema file
	content, err := fs.ReadFile(fsys, path)
	if err != nil {
		return MissingSchemaError{}
	}

	// Create URI for this schema (version + filename)
	uri := fmt.Sprintf("acdc-schema:///%s/%s", version, filepath.Base(path))
	schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(content))
	if err != nil {
		return err
	}
	// Add to compiler
	if err := vc.compiler.AddResource(uri, schema); err != nil {
		return err
	}

	return nil
}

// chooseFile picks schema filename from the credential event map
func chooseFile(ked map[string]interface{}) (string, error) {
	a, ok := ked["a"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("missing \"a\" section")
	}
	switch {
	case a["engagementContextRole"] != nil && a["AID"] != nil:
		return "ecr-authorization-vlei-credential.json", nil

	case a["engagementContextRole"] != nil:
		return "legal-entity-engagement-context-role-vLEI-credential.json", nil

	case a["officialRole"] != nil && a["AID"] != nil:
		return "oor-authorization-vlei-credential.json", nil

	case a["officialRole"] != nil:
		return "legal-entity-official-organizational-role-vLEI-credential.json", nil

	case a["LEI"] != nil && ked["e"] != nil:
		return "legal-entity-vLEI-credential.json", nil

	case a["LEI"] != nil:
		return "qualified-vLEI-issuer-vLEI-credential.json", nil

	case a["f"] != nil:
		return "verifiable-ixbrl-report-attestation.json", nil

	default:
		return "", fmt.Errorf("unrecognised credential type")
	}
}

// ValidateCredential tries all schema dirs until one validates the credential.
// Returns nil if successful, or an aggregate error listing attempts.
// ValidateCredential uses pre-loaded compilers to validate schemas
func (v *Validator) ValidateCredential(ked map[string]interface{}) error {
	// Skip non-ACDC events
	if ver, _ := ked["v"].(string); !strings.HasPrefix(ver, "ACDC") {
		return nil
	}

	fname, err := chooseFile(ked)
	if err != nil {
		return err
	}

	var attempts []string
	for _, vc := range v.versions {
		// Create URI for the schema file in this version
		uri := fmt.Sprintf("acdc-schema:///%s/%s", vc.name, fname)
		path := filepath.Join(vc.name, fname)
		err = vc.updateCompiler(v.root, vc.name, path)
		if err != nil {
			if _, ok := err.(MissingSchemaError); !ok {
				attempts = append(attempts, fmt.Sprintf("%s (compile error: %v)", uri, err))
			}
			continue
		}
		sch, err := vc.compiler.Compile(uri)
		if err != nil {
			attempts = append(attempts,
				fmt.Sprintf("%s (compiled error: %v)", uri, err))
			continue
		}

		// Validate against schema
		if err = sch.Validate(ked); err == nil {
			return nil // Success
		}
		attempts = append(attempts,
			fmt.Sprintf("%s (validation error: %v)", uri, err))

	}

	if len(attempts) == 0 {
		return fmt.Errorf("no schema versions available for validation")
	}

	return fmt.Errorf("no schema accepted credential; tried:\n  %s",
		strings.Join(attempts, "\n  "))
}
