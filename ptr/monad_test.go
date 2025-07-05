package ptr_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-ptr-tools/ptr"
)

var errNegativeNumber = errors.New("negative number")

// Helper functions for testing.
func add1(x int) int { return x + 1 }

func addCtx(ctx context.Context, number int) int {
	select {
	default:
	case <-ctx.Done():
		return 0 // return zero-value if context is done
	}

	return number + 1
}

func addErr(x int) (int, error) {
	if x < 0 {
		return 0, errNegativeNumber
	}

	return x + 1, nil
}

func addCtxErr(ctx context.Context, x int) (int, error) {
	select {
	default:
	case <-ctx.Done():
		return 0, fmt.Errorf("context done, no numbers were incremented: %w", ctx.Err())
	}

	if x < 0 {
		return 0, errNegativeNumber
	}

	return x + 1, nil
}

func TestApply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *int
		expected *int
	}{
		{
			name:     "normal case",
			input:    ptr.New(5),
			expected: ptr.New(6),
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ptr.Apply(tt.input, add1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestApplyCtx(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		name     string
		input    *int
		expected *int
	}{
		{
			name:     "normal case",
			input:    ptr.New(5),
			expected: ptr.New(6),
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ptr.ApplyCtx(ctx, tt.input, addCtx)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestApplyErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       *int
		expected    *int
		expectError error
	}{
		{
			name:        "normal case",
			input:       ptr.New(5),
			expected:    ptr.New(6),
			expectError: nil,
		},
		{
			name:        "error case",
			input:       ptr.New(-1),
			expected:    nil,
			expectError: errNegativeNumber,
		},
		{
			name:        "nil input",
			input:       nil,
			expected:    nil,
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := ptr.ApplyErr(tt.input, addErr)
			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
