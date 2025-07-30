package cesr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const dummySchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Legal Entity vLEI Credential",
  "type": "object",
  "required": ["v", "a", "e"],
  "properties": {
    "v": {
      "type": "string"
    },
    "a": {
      "type": "object",
      "required": ["LEI"],
      "properties": {
        "LEI": {
          "type": "string"
        }
      },
      "additionalProperties": false
    },
    "e": {
      "type": "object"
    }
  },
  "additionalProperties": false
}`

func writeSchemaFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write schema file: %v", err)
	}
}

func TestValidateCredential(t *testing.T) {
	tests := []struct {
		name          string
		credential    map[string]interface{}
		setupSchemas  func(t *testing.T, root string)
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid credential succeeds",
			credential: map[string]interface{}{
				"v": "ACDC000000_",
				"a": map[string]interface{}{
					"LEI": "1234567890ABCDEF",
				},
				"e": map[string]interface{}{
					"something": true,
				},
			},
			setupSchemas: func(t *testing.T, root string) {
				dir := filepath.Join(root, "2023")
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
				writeSchemaFile(t, dir, "legal-entity-vLEI-credential.json", dummySchema)
			},
			expectError: false,
		},
		{
			name: "Non-ACDC version is ignored",
			credential: map[string]interface{}{
				"v": "KERI100000_",
				"a": map[string]interface{}{
					"LEI": "xyz",
				},
			},
			setupSchemas: func(t *testing.T, root string) {},
			expectError:  false,
		},
		{
			name: "Unrecognized credential type",
			credential: map[string]interface{}{
				"v": "ACDC/1.0",
				"a": map[string]interface{}{
					"unknownField": true,
				},
			},
			setupSchemas:  func(t *testing.T, root string) {},
			expectError:   true,
			errorContains: "unrecognised credential type",
		},
		{
			name: "Schema compile error",
			credential: map[string]interface{}{
				"v": "ACDC/1.0",
				"a": map[string]interface{}{
					"LEI": "1234567890ABCDEF",
				},
				"e": map[string]interface{}{
					"something": true,
				},
			},
			setupSchemas: func(t *testing.T, root string) {
				dir := filepath.Join(root, "2023")
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
				writeSchemaFile(t, dir, "legal-entity-vLEI-credential.json", `{ invalid json `)
			},
			expectError:   true,
			errorContains: "compile error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpRoot := t.TempDir()
			tt.setupSchemas(t, tmpRoot)

			v := NewValidator(tmpRoot)
			err := v.ValidateCredential(tt.credential)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Fatalf("expected error to contain %q, got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected success but got error: %v", err)
				}
			}
		})
	}
}
