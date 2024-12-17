package opt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
)

func testMakeIntHelper(exp int, action func(*int)) func(*testing.T, *int) {
	return func(t *testing.T, v *int) {
		t.Helper()
		require.Equal(t, &exp, v)
		action(v)
	}
}

func testMakeIntHelperEx(exp int, err error, action func(*int)) func(*testing.T, *int) error {
	return func(t *testing.T, v *int) error {
		t.Helper()
		require.Equal(t, &exp, v)
		action(v)

		return err
	}
}

func testMakeSliceHelper(exp, signal int) func(*testing.T, int, *[]int) {
	return func(t *testing.T, val int, s *[]int) {
		t.Helper()
		require.Equal(t, exp, val)

		*s = append(*s, signal)
	}
}

func TestTrigger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		optional opt.Opt[int]
		funcs    []func(*testing.T, int, *[]int)
		expected []int
	}{
		{
			"single functions",
			opt.Of(42),
			[]func(*testing.T, int, *[]int){testMakeSliceHelper(42, 1)},
			[]int{1},
		},
		{
			"multiple functions",
			opt.Of(55),
			[]func(*testing.T, int, *[]int){
				testMakeSliceHelper(55, 10),
				testMakeSliceHelper(55, 20),
			},
			[]int{10, 20},
		},
		{
			"empty",
			opt.Empty[int](),
			[]func(*testing.T, int, *[]int){func(t *testing.T, _ int, _ *[]int) { t.Helper(); t.Fail() }},
			[]int{},
		},
		{
			"no functions",
			opt.Of(42),
			[]func(*testing.T, int, *[]int){},
			[]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				result []int
				funcs  []func(int)
			)

			for _, f := range tt.funcs {
				funcs = append(funcs, func(v int) { f(t, v, &result) })
			}

			tt.optional.Trigger(funcs...)
			require.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestApply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		optional opt.Opt[int]
		funcs    []func(*testing.T, *int)
		expected *int
	}{
		{
			"single function",
			opt.Of(42),
			[]func(*testing.T, *int){testMakeIntHelper(42, func(v *int) { *v++ })},
			ptr.Of(43),
		},
		{
			"multiple functions",
			opt.Of(55),
			[]func(*testing.T, *int){
				testMakeIntHelper(55, func(v *int) { *v += 10 }),
				testMakeIntHelper(65, func(v *int) { *v *= 2 }),
			},
			ptr.Of(130),
		},
		{
			"empty",
			opt.Empty[int](),
			[]func(*testing.T, *int){func(t *testing.T, _ *int) { t.Helper(); t.Fail() }},
			nil,
		},
		{
			"no functions",
			opt.Of(42),
			[]func(*testing.T, *int){},
			ptr.Of(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var funcs []func(*int)
			for _, f := range tt.funcs {
				funcs = append(funcs, func(v *int) { f(t, v) })
			}

			tt.optional.Apply(funcs...)
			require.Equal(t, tt.expected, tt.optional.Ptr())
		})
	}
}

func TestApplyEx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		optional opt.Opt[int]
		funcs    []func(*testing.T, *int) error
		expected *int
		err      error
	}{
		{
			"single function",
			opt.Of(42),
			[]func(*testing.T, *int) error{testMakeIntHelperEx(42, nil, func(v *int) { *v++ })},
			ptr.Of(43),
			nil,
		},
		{
			"multiple functions",
			opt.Of(55),
			[]func(*testing.T, *int) error{
				testMakeIntHelperEx(55, nil, func(v *int) { *v += 10 }),
				testMakeIntHelperEx(65, nil, func(v *int) { *v *= 2 }),
			},
			ptr.Of(130),
			nil,
		},
		{
			"error",
			opt.Of(55),
			[]func(*testing.T, *int) error{
				testMakeIntHelperEx(55, nil, func(v *int) { *v += 10 }),
				testMakeIntHelperEx(65, errTest, func(_ *int) {}),
				func(t *testing.T, _ *int) error { t.Helper(); t.Fail(); return nil },
			},
			ptr.Of(65),
			errTest,
		},
		{
			"empty",
			opt.Empty[int](),
			[]func(*testing.T, *int) error{
				func(t *testing.T, _ *int) error { t.Helper(); t.Fail(); return nil },
			},
			nil,
			nil,
		},
		{
			"no functions",
			opt.Of(42),
			[]func(*testing.T, *int) error{},
			ptr.Of(42),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var funcs []func(*int) error
			for _, f := range tt.funcs {
				funcs = append(funcs, func(v *int) error { return f(t, v) })
			}

			err := tt.optional.ApplyEx(funcs...)
			require.ErrorIs(t, err, tt.err)

			if err == nil && tt.err == nil {
				require.Equal(t, tt.expected, tt.optional.Ptr())
			}
		})
	}
}

func TestElse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		optional opt.Opt[int]
		val      int
		expected int
	}{
		{
			"value 42",
			opt.Of(42),
			100,
			42,
		},
		{
			"empty",
			opt.Empty[int](),
			100,
			100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.optional.Else(tt.val)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestElseGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		optional opt.Opt[int]
		f        func(*testing.T) int
		expected int
	}{
		{
			"value 42",
			opt.Of(42),
			func(t *testing.T) int { t.Helper(); t.Fail(); return 0 },
			42,
		},
		{
			"empty",
			opt.Empty[int](),
			func(_ *testing.T) int { return 100 },
			100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tt.optional.ElseGet(func() int { return tt.f(t) })
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestElseGetEx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		optional opt.Opt[int]
		f        func(*testing.T) (int, error)
		expected int
		err      error
	}{
		{
			"value 42",
			opt.Of(42),
			func(t *testing.T) (int, error) { t.Helper(); t.Fail(); return 0, nil },
			42,
			nil,
		},
		{
			"empty",
			opt.Empty[int](),
			func(_ *testing.T) (int, error) { return 100, nil },
			100,
			nil,
		},
		{
			"error",
			opt.Empty[int](),
			func(_ *testing.T) (int, error) { return 0, errTest },
			0,
			errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := tt.optional.ElseGetEx(func() (int, error) { return tt.f(t) })
			require.ErrorIs(t, err, tt.err)

			if err == nil && tt.err == nil {
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
