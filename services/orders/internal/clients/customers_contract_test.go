package clients

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// isActiveResponse is the body customers-service answers POST /v1/customers/is_active with.
type isActiveResponse struct {
	Active bool `json:"active"`
}

// Every golden customers-service response decodes with no unknown fields.
func TestCustomersContractFixtures(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "contracts", "customers", "*.json"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, f := range files {
		t.Run(filepath.Base(f), func(t *testing.T) {
			raw, err := os.ReadFile(f)
			require.NoError(t, err)

			dec := json.NewDecoder(bytes.NewReader(raw))
			dec.DisallowUnknownFields()
			var r isActiveResponse
			require.NoError(t, dec.Decode(&r))
			require.True(t, r.Active)
		})
	}
}
