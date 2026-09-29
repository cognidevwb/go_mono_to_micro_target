package clients

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	contract "github.com/acme/shop/services/orders/internal/contracts/catalog"
)

// Every golden catalog-service response decodes into orders's contract copy
// with no unknown fields and satisfies the invariants the types cannot express.
func TestCatalogContractFixtures(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "contracts", "catalog", "*.json"))
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, f := range files {
		t.Run(filepath.Base(f), func(t *testing.T) {
			raw, err := os.ReadFile(f)
			require.NoError(t, err)

			dec := json.NewDecoder(bytes.NewReader(raw))
			dec.DisallowUnknownFields()
			var p contract.Product
			require.NoError(t, dec.Decode(&p))

			require.NotZero(t, p.ID)
			require.NotEmpty(t, p.SKU)
			require.NotEmpty(t, p.Name)
			require.Greater(t, p.Price, 0.0)
		})
	}
}
