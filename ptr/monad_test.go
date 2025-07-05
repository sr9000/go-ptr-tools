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

func addCtxErr(ctx context.Context, number int) (int, error) {
	select {
	default:
	case <-ctx.Done():
		return 0, fmt.Errorf("context done, no numbers were incremented: %w", ctx.Err())
	}

	if number < 0 {
		return 0, errNegativeNumber
	}

	return number + 1, nil
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

func TestApplyCtxErr(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
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

			result, err := ptr.ApplyCtxErr(ctx, tt.input, addCtxErr)
			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestApplyVoid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *int
		expected bool // tracks if function was called
	}{
		{
			name:     "normal case",
			input:    ptr.New(5),
			expected: true,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wasCalled := false

			ptr.ApplyVoid(tt.input, func(x int) {
				wasCalled = true

				require.Equal(t, *tt.input, x)
			})
			require.Equal(t, tt.expected, wasCalled)
		})
	}
}

func TestApplyVoidCtx(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		name        string
		input       *int
		expCalled   bool // tracks if function was called
		expAnalyzed bool // tracks if number were analyzed
	}{
		{
			name:        "normal case",
			input:       ptr.New(5),
			expCalled:   true,
			expAnalyzed: true,
		},
		{
			name:        "nil input",
			input:       nil,
			expCalled:   false,
			expAnalyzed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wasCalled := false
			wasAnalyzed := false

			ptr.ApplyVoidCtx(ctx, tt.input, func(ctx context.Context, number int) {
				wasCalled = true

				select {
				default:
				case <-ctx.Done():
					return // return if context is done
				}

				wasAnalyzed = true

				require.Equal(t, *tt.input, number)
			})

			require.Equal(t, tt.expCalled, wasCalled)
			require.Equal(t, tt.expAnalyzed, wasAnalyzed)
		})
	}
}

func TestApplyVoidErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       *int
		expectError error
		expected    bool // tracks if function was called
	}{
		{
			name:        "normal case",
			input:       ptr.New(5),
			expectError: nil,
			expected:    true,
		},
		{
			name:        "error case",
			input:       ptr.New(-1),
			expectError: errNegativeNumber,
			expected:    true,
		},
		{
			name:        "nil input",
			input:       nil,
			expectError: nil,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wasCalled := false
			err := ptr.ApplyVoidErr(tt.input, func(x int) error {
				wasCalled = true

				require.Equal(t, *tt.input, x)

				if x < 0 {
					return errNegativeNumber
				}

				return nil
			})

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expected, wasCalled)
		})
	}
}

func TestApplyVoidCtxErr(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	tests := []struct {
		name        string
		input       *int
		expectError error
		expCalled   bool // tracks if function was called
		expAnalyzed bool // tracks if function was analyzed
	}{
		{
			name:        "normal case",
			input:       ptr.New(5),
			expectError: nil,
			expCalled:   true,
			expAnalyzed: true,
		},
		{
			name:        "error case",
			input:       ptr.New(-1),
			expectError: errNegativeNumber,
			expCalled:   true,
			expAnalyzed: true,
		},
		{
			name:        "nil input",
			input:       nil,
			expectError: nil,
			expCalled:   false,
			expAnalyzed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wasCalled := false
			wasAnalyzed := false
			err := ptr.ApplyVoidCtxErr(ctx, tt.input, func(ctx context.Context, number int) error {
				wasCalled = true

				select {
				default:
				case <-ctx.Done():
					return fmt.Errorf("context done, no numbers were analyzed: %w", ctx.Err())
				}

				wasAnalyzed = true

				require.Equal(t, *tt.input, number)

				if number < 0 {
					return errNegativeNumber
				}

				return nil
			})

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expCalled, wasCalled)
			require.Equal(t, tt.expAnalyzed, wasAnalyzed)
		})
	}
}

// Helper functions for testing Apply9.
func add9(a, b, c, d, e, f, g, h, i int) int {
	return a + b + c + d + e + f + g + h + i
}

func TestApply9(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputs   [9]*int
		expected *int
	}{
		{
			name: "normal case",
			inputs: [9]*int{
				ptr.New(1), ptr.New(2), ptr.New(3),
				ptr.New(4), ptr.New(5), ptr.New(6),
				ptr.New(7), ptr.New(8), ptr.New(9),
			},
			expected: ptr.New(45),
		},
		{
			name:     "nil input",
			inputs:   [9]*int{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ptr.Apply9(
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				add9,
			)
			require.Equal(t, tt.expected, result)
		})
	}
}
