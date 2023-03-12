package isitcg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadRules(t *testing.T) {
	got, err := LoadRules("../../ingredientrules.json")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Len(t, got, 24)
	require.Len(t, got[15].Ingredients, 86)
}
