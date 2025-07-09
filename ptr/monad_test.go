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

// Helper function to create cancelled context.
func cancelledCtx(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	return ctx
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
			input:    ptr.Of(5),
			expected: ptr.Of(6),
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

			res := ptr.Apply(tt.input, add1)
			require.Equal(t, tt.expected, res)
		})
	}
}

func TestApplyCtx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		input        *int
		expected     *int
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			input:        ptr.Of(5),
			expected:     ptr.Of(6),
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			input:        ptr.Of(5),
			expected:     ptr.Of(0), // zero value when context cancelled
		},
		{
			name:         "nil input",
			ctxCancelled: false,
			input:        nil,
			expected:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			res := ptr.ApplyCtx(ctx, tt.input, addCtx)
			require.Equal(t, tt.expected, res)
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
			input:       ptr.Of(5),
			expected:    ptr.Of(6),
			expectError: nil,
		},
		{
			name:        "error case",
			input:       ptr.Of(-1),
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

			res, err := ptr.ApplyErr(tt.input, addErr)
			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, res)
			}
		})
	}
}

func TestApplyCtxErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		input        *int
		expected     *int
		expectError  error
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			input:        ptr.Of(5),
			expected:     ptr.Of(6),
			expectError:  nil,
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			input:        ptr.Of(5),
			expected:     nil,
			expectError:  context.Canceled,
		},
		{
			name:         "error case",
			ctxCancelled: false,
			input:        ptr.Of(-1),
			expected:     nil,
			expectError:  errNegativeNumber,
		},
		{
			name:         "nil input",
			ctxCancelled: false,
			input:        nil,
			expected:     nil,
			expectError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			res, err := ptr.ApplyCtxErr(ctx, tt.input, addCtxErr)

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, res)
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
			input:    ptr.Of(5),
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
			input:       ptr.Of(5),
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
			input:       ptr.Of(5),
			expectError: nil,
			expected:    true,
		},
		{
			name:        "error case",
			input:       ptr.Of(-1),
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
			input:       ptr.Of(5),
			expectError: nil,
			expCalled:   true,
			expAnalyzed: true,
		},
		{
			name:        "error case",
			input:       ptr.Of(-1),
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

func ptrsArray(xs []int) []*int {
	// Helper function to convert slice of ints to slice of *int
	ints := make([]*int, len(xs))

	for i, x := range xs {
		if x != 0 {
			ints[i] = &x
		}
	}

	return ints
}

func TestApply9Void(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputs   []*int
		expected bool // tracks if function was called
	}{
		{
			name:     "normal case",
			inputs:   ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected: true,
		},
		{
			name:     "one nil input",
			inputs:   ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected: false,
		},
		{
			name:     "all nil inputs",
			inputs:   ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wasCalled := false

			ptr.Apply9Void(
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(num1, num2, num3, num4, num5, num6, num7, num8, num9 int) {
					wasCalled = true

					require.Equal(t, []*int{&num1, &num2, &num3, &num4, &num5, &num6, &num7, &num8, &num9}, tt.inputs)
				},
			)
			require.Equal(t, tt.expected, wasCalled)
		})
	}
}

func TestApply9VoidCtx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		inputs       []*int
		expCalled    bool // tracks if function was called
		expAnalyzed  bool // tracks if numbers were analyzed
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expCalled:    true,
			expAnalyzed:  true,
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expCalled:    true,
			expAnalyzed:  false,
		},
		{
			name:         "one nil input",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expCalled:    false,
			expAnalyzed:  false,
		},
		{
			name:         "all nil inputs",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expCalled:    false,
			expAnalyzed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			wasCalled := false
			wasAnalyzed := false

			ptr.Apply9VoidCtx(
				ctx,
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(ctx context.Context, num1, num2, num3, num4, num5, num6, num7, num8, num9 int) {
					wasCalled = true

					select {
					default:
					case <-ctx.Done():
						return
					}

					wasAnalyzed = true

					require.Equal(t, []*int{&num1, &num2, &num3, &num4, &num5, &num6, &num7, &num8, &num9}, tt.inputs)
				},
			)

			require.Equal(t, tt.expCalled, wasCalled)
			require.Equal(t, tt.expAnalyzed, wasAnalyzed)
		})
	}
}

func TestApply9VoidErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputs      []*int
		expectError error
		expCalled   bool // tracks if function was called
	}{
		{
			name:        "normal case",
			inputs:      ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expectError: nil,
			expCalled:   true,
		},
		{
			name:        "error case",
			inputs:      ptrsArray([]int{1, 2, 3, 4, -5, 6, 7, 8, 9}),
			expectError: errNegativeNumber,
			expCalled:   true,
		},
		{
			name:        "one nil input",
			inputs:      ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expectError: nil,
			expCalled:   false,
		},
		{
			name:        "all nil inputs",
			inputs:      ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expectError: nil,
			expCalled:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			wasCalled := false
			err := ptr.Apply9VoidErr(
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(num1, num2, num3, num4, num5, num6, num7, num8, num9 int) error {
					wasCalled = true

					require.Equal(t, []*int{&num1, &num2, &num3, &num4, &num5, &num6, &num7, &num8, &num9}, tt.inputs)

					if min(num1, num2, num3, num4, num5, num6, num7, num8, num9) < 0 {
						return errNegativeNumber
					}

					return nil
				},
			)

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expCalled, wasCalled)
		})
	}
}

func TestApply9VoidCtxErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		inputs       []*int
		expectError  error
		expCalled    bool // tracks if function was called
		expAnalyzed  bool // tracks if numbers were analyzed
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expectError:  nil,
			expCalled:    true,
			expAnalyzed:  true,
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expectError:  context.Canceled,
			expCalled:    true,
			expAnalyzed:  false,
		},
		{
			name:         "error case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, -5, 6, 7, 8, 9}),
			expectError:  errNegativeNumber,
			expCalled:    true,
			expAnalyzed:  true,
		},
		{
			name:         "one nil input",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expectError:  nil,
			expCalled:    false,
			expAnalyzed:  false,
		},
		{
			name:         "all nil inputs",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expectError:  nil,
			expCalled:    false,
			expAnalyzed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			wasCalled := false
			wasAnalyzed := false
			err := ptr.Apply9VoidCtxErr(
				ctx,
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(ctx context.Context, num1, num2, num3, num4, num5, num6, num7, num8, num9 int) error {
					wasCalled = true

					select {
					default:
					case <-ctx.Done():
						return ctx.Err()
					}

					wasAnalyzed = true

					require.Equal(t, []*int{&num1, &num2, &num3, &num4, &num5, &num6, &num7, &num8, &num9}, tt.inputs)

					if min(num1, num2, num3, num4, num5, num6, num7, num8, num9) < 0 {
						return errNegativeNumber
					}

					return nil
				},
			)

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

func TestApply9(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputs   []*int
		expected *int
	}{
		{
			name:     "normal case",
			inputs:   ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected: ptr.Of(45), // sum of all numbers
		},
		{
			name:     "one nil input",
			inputs:   ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected: nil,
		},
		{
			name:     "all nil inputs",
			inputs:   ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := ptr.Apply9(
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(num1, num2, num3, num4, num5, num6, num7, num8, num9 int) int {
					return num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9
				})

			require.Equal(t, tt.expected, res)
		})
	}
}

func TestApply9Ctx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		inputs       []*int
		expected     *int
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected:     ptr.Of(45), // sum of all numbers
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected:     ptr.Of(0), // zero value when context cancelled
		},
		{
			name:         "one nil input",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected:     nil,
		},
		{
			name:         "all nil inputs",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			res := ptr.Apply9Ctx(
				ctx,
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(ctx context.Context, num1, num2, num3, num4, num5, num6, num7, num8, num9 int) int {
					select {
					default:
					case <-ctx.Done():
						return 0
					}

					return num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9
				})

			require.Equal(t, tt.expected, res)
		})
	}
}

func TestApply9Err(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputs      []*int
		expected    *int
		expectError error
	}{
		{
			name:        "normal case",
			inputs:      ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected:    ptr.Of(45), // sum of all numbers
			expectError: nil,
		},
		{
			name:        "error case",
			inputs:      ptrsArray([]int{1, 2, 3, 4, -5, 6, 7, 8, 9}),
			expected:    nil,
			expectError: errNegativeNumber,
		},
		{
			name:        "one nil input",
			inputs:      ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected:    nil,
			expectError: nil,
		},
		{
			name:        "all nil inputs",
			inputs:      ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected:    nil,
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := ptr.Apply9Err(
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(num1, num2, num3, num4, num5, num6, num7, num8, num9 int) (int, error) {
					if min(num1, num2, num3, num4, num5, num6, num7, num8, num9) < 0 {
						return 0, errNegativeNumber
					}

					return num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9, nil
				})

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, res)
			}
		})
	}
}

func TestApply9CtxErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		inputs       []*int
		expected     *int
		expectError  error
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected:     ptr.Of(45), // sum of all numbers
			expectError:  nil,
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected:     nil,
			expectError:  context.Canceled,
		},
		{
			name:         "error case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, -5, 6, 7, 8, 9}),
			expected:     nil,
			expectError:  errNegativeNumber,
		},
		{
			name:         "one nil input",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected:     nil,
			expectError:  nil,
		},
		{
			name:         "all nil inputs",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected:     nil,
			expectError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			res, err := ptr.Apply9CtxErr(
				ctx,
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(ctx context.Context, num1, num2, num3, num4, num5, num6, num7, num8, num9 int) (int, error) {
					select {
					default:
					case <-ctx.Done():
						return 0, ctx.Err()
					}

					if min(num1, num2, num3, num4, num5, num6, num7, num8, num9) < 0 {
						return 0, errNegativeNumber
					}

					return num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9, nil
				})

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, res)
			}
		})
	}
}

type result struct {
	sum, product, min, max, avg int
}

func requireAllNil(t *testing.T, sum *int, product *int, mn *int, mx *int, avg *int) {
	t.Helper()

	require.Nil(t, sum)
	require.Nil(t, product)
	require.Nil(t, mn)
	require.Nil(t, mx)
	require.Nil(t, avg)
}

func requireExpected(t *testing.T, expected result, sum *int, product *int, mn *int, mx *int, avg *int) {
	t.Helper()

	require.Equal(t, expected.sum, *sum)
	require.Equal(t, expected.product, *product)
	require.Equal(t, expected.min, *mn)
	require.Equal(t, expected.max, *mx)
	require.Equal(t, expected.avg, *avg)
}

func TestApply95(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		inputs   []*int
		expected *result
	}{
		{
			name:   "normal case",
			inputs: ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected: &result{
				sum:     45,
				product: 362880,
				min:     1,
				max:     9,
				avg:     5,
			},
		},
		{
			name:     "one nil input",
			inputs:   ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected: nil,
		},
		{
			name:     "all nil inputs",
			inputs:   ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sum, product, minimum, maximum, avg := ptr.Apply95(
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(num1, num2, num3, num4, num5, num6, num7, num8, num9 int) (sum, product, mn, mx, avg int) {
					sum = num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9
					product = num1 * num2 * num3 * num4 * num5 * num6 * num7 * num8 * num9
					mn = min(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					mx = max(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					avg = sum / 9

					return
				})

			if tt.expected == nil {
				requireAllNil(t, sum, product, minimum, maximum, avg)
			} else {
				requireExpected(t, *tt.expected, sum, product, minimum, maximum, avg)
			}
		})
	}
}

func TestApply95Ctx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		inputs       []*int
		expected     *result
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected: &result{
				sum:     45,
				product: 362880,
				min:     1,
				max:     9,
				avg:     5,
			},
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected: &result{
				sum:     0,
				product: 0,
				min:     0,
				max:     0,
				avg:     0,
			},
		},
		{
			name:         "one nil input",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected:     nil,
		},
		{
			name:         "all nil inputs",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			sum, product, minimum, maximum, avg := ptr.Apply95Ctx(
				ctx,
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(ctx context.Context, num1, num2, num3, num4, num5, num6, num7, num8, num9 int,
				) (sum, product, mn, mx, avg int) {
					select {
					default:
					case <-ctx.Done():
						return 0, 0, 0, 0, 0
					}

					sum = num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9
					product = num1 * num2 * num3 * num4 * num5 * num6 * num7 * num8 * num9
					mn = min(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					mx = max(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					avg = sum / 9

					return
				})

			if tt.expected == nil {
				requireAllNil(t, sum, product, minimum, maximum, avg)
			} else {
				requireExpected(t, *tt.expected, sum, product, minimum, maximum, avg)
			}
		})
	}
}

func TestApply95Err(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputs      []*int
		expected    *result
		expectError error
	}{
		{
			name:   "normal case",
			inputs: ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected: &result{
				sum:     45,
				product: 362880,
				min:     1,
				max:     9,
				avg:     5,
			},
			expectError: nil,
		},
		{
			name:        "error case",
			inputs:      ptrsArray([]int{1, 2, 3, 4, -5, 6, 7, 8, 9}),
			expected:    nil,
			expectError: errNegativeNumber,
		},
		{
			name:        "one nil input",
			inputs:      ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected:    nil,
			expectError: nil,
		},
		{
			name:        "all nil inputs",
			inputs:      ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected:    nil,
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sum, product, minimum, maximum, avg, err := ptr.Apply95Err(
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(num1, num2, num3, num4, num5, num6, num7, num8, num9 int) (sum, product, mn, mx, avg int, err error) {
					if min(num1, num2, num3, num4, num5, num6, num7, num8, num9) < 0 {
						return 0, 0, 0, 0, 0, errNegativeNumber
					}

					sum = num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9
					product = num1 * num2 * num3 * num4 * num5 * num6 * num7 * num8 * num9
					mn = min(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					mx = max(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					avg = sum / 9

					return
				})

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)

				if tt.expected == nil {
					requireAllNil(t, sum, product, minimum, maximum, avg)
				} else {
					requireExpected(t, *tt.expected, sum, product, minimum, maximum, avg)
				}
			}
		})
	}
}

func TestApply95CtxErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		ctxCancelled bool
		inputs       []*int
		expected     *result
		expectError  error
	}{
		{
			name:         "normal case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected: &result{
				sum:     45,
				product: 362880,
				min:     1,
				max:     9,
				avg:     5,
			},
			expectError: nil,
		},
		{
			name:         "cancelled context",
			ctxCancelled: true,
			inputs:       ptrsArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			expected:     nil,
			expectError:  context.Canceled,
		},
		{
			name:         "error case",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 2, 3, 4, -5, 6, 7, 8, 9}),
			expected:     nil,
			expectError:  errNegativeNumber,
		},
		{
			name:         "one nil input",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{1, 0, 3, 4, 5, 6, 7, 8, 9}),
			expected:     nil,
			expectError:  nil,
		},
		{
			name:         "all nil inputs",
			ctxCancelled: false,
			inputs:       ptrsArray([]int{0, 0, 0, 0, 0, 0, 0, 0, 0}),
			expected:     nil,
			expectError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			if tt.ctxCancelled {
				ctx = cancelledCtx(t)
			}

			sum, product, minimum, maximum, avg, err := ptr.Apply95CtxErr(
				ctx,
				tt.inputs[0], tt.inputs[1], tt.inputs[2],
				tt.inputs[3], tt.inputs[4], tt.inputs[5],
				tt.inputs[6], tt.inputs[7], tt.inputs[8],
				func(ctx context.Context, num1, num2, num3, num4, num5, num6, num7, num8, num9 int,
				) (sum, product, mn, mx, avg int, err error) {
					select {
					default:
					case <-ctx.Done():
						return 0, 0, 0, 0, 0, ctx.Err()
					}

					if min(num1, num2, num3, num4, num5, num6, num7, num8, num9) < 0 {
						return 0, 0, 0, 0, 0, errNegativeNumber
					}

					sum = num1 + num2 + num3 + num4 + num5 + num6 + num7 + num8 + num9
					product = num1 * num2 * num3 * num4 * num5 * num6 * num7 * num8 * num9
					mn = min(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					mx = max(num1, num2, num3, num4, num5, num6, num7, num8, num9)
					avg = sum / 9

					return
				})

			if tt.expectError != nil {
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)

				if tt.expected == nil {
					requireAllNil(t, sum, product, minimum, maximum, avg)
				} else {
					requireExpected(t, *tt.expected, sum, product, minimum, maximum, avg)
				}
			}
		})
	}
}
