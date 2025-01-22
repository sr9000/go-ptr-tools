package examples_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBug(t *testing.T) {
	t.Parallel()

	var (
		a, b int
		c    int
	)

	a, b = 1, 2
	c = a + b

	require.Equal(t, 3, c)
}
