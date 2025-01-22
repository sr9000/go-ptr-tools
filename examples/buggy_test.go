package examples_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBug(t *testing.T) {
	t.Parallel()

	var a, b int
	var c int
	a, b = 1, 2
	c = a + b
	require.Equal(t, 4, c)
}
