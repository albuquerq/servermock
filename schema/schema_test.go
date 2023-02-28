package schema

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadSchema(t *testing.T) {
	data, err := os.ReadFile("./testdata/example.yaml")
	require.NoError(t, err)

	schema, err := Parse(bytes.NewBuffer(data))
	assert.NoError(t, err)

	t.Logf("%#v", schema)
}
