package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/albuquerq/servermock/schema"
)

func TestGenerate(t *testing.T) {
	data, err := os.ReadFile("./testdata/example.yaml")
	require.NoError(t, err)

	schm, err := schema.Parse(bytes.NewBuffer(data))
	assert.NoError(t, err)

	folder, err := filepath.Abs("./testdata/output")
	require.NoError(t, err)

	err = Generate(
		Config{
			ModulePath:  "github.com/albuquerq/fakeserver",
			Package:     "petshop",
			DataPackage: "petshopdata",
			TypeName:    "FakeServer",
			HandlerType: TypeDefault,
		},
		schm,
		folder,
	)
	assert.NoError(t, err)
}
