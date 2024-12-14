package opt_test

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
	"github.com/sr9000/go-noptr/noptr/val"
)

func BenchmarkCastTo(b *testing.B) {
	bar := &testBar{}
	nul := (*testBar)(nil)
	str := "string"

	for range b.N / 4 {
		_ = opt.CastTo[testFooer](bar)
		_ = opt.CastTo[testFooer](nul)
		_ = opt.CastTo[testFooer](nil)
		_ = opt.CastTo[testFooer](str)
	}
}

func TestCastTo_Int(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		value    any
		expected *int
	}{
		{"valid cast", 42, ptr.Of(42)},
		{"invalid cast", "string", nil},
		{"nil value", nil, nil},
		{"nil pointer", ptr.Nil[int](), nil},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.CastTo[int](cs.value)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}

func TestCastTo_SliceOfString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		value    any
		expected *[]string
	}{
		{"valid cast", []string{"foo", "bar"}, ptr.Of([]string{"foo", "bar"})},
		{"invalid cast", "string", nil},
		{"nil value", nil, nil},
		{"nil slice ptr", ptr.Nil[[]string](), nil},
		{"zero slice", val.Zero[[]string](), new([]string)},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.CastTo[[]string](cs.value)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}

func TestCoalesce(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		opts     []opt.Opt[int]
		expected opt.Opt[int]
	}{
		{
			"all empty",
			[]opt.Opt[int]{opt.Empty[int](), opt.Empty[int](), opt.Empty[int]()},
			opt.Empty[int](),
		},
		{
			"first not empty",
			[]opt.Opt[int]{opt.Of(1), opt.Empty[int](), opt.Empty[int]()},
			opt.Of(1),
		},
		{
			"second not empty",
			[]opt.Opt[int]{opt.Empty[int](), opt.Of(2), opt.Empty[int]()},
			opt.Of(2),
		},
		{
			"third not empty",
			[]opt.Opt[int]{opt.Empty[int](), opt.Empty[int](), opt.Of(3)},
			opt.Of(3),
		},
		{
			"all not empty",
			[]opt.Opt[int]{opt.Of(1), opt.Of(2), opt.Of(3)},
			opt.Of(1),
		},
		{
			"single not empty",
			[]opt.Opt[int]{opt.Of(1)},
			opt.Of(1),
		},
		{
			"single empty",
			[]opt.Opt[int]{opt.Empty[int]()},
			opt.Empty[int](),
		},
		{
			"no optionals",
			[]opt.Opt[int]{},
			opt.Empty[int](),
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.Coalesce(cs.opts...)
			require.Equal(t, cs.expected, result)
		})
	}
}

func TestFMap_SingleTransform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		opt      opt.Opt[int]
		f        func(t *testing.T, i int) string
		expected opt.Opt[string]
	}{
		{
			"value 42",
			opt.Of(42),
			func(t *testing.T, i int) string {
				t.Helper()
				require.Equal(t, 42, i)

				return fmt.Sprintf("Value: %d", i)
			},
			opt.Of("Value: 42"),
		},
		{
			"empty",
			opt.Empty[int](),
			func(t *testing.T, _ int) string { t.Helper(); t.Fail(); return "" },
			opt.Empty[string](),
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.FMap(cs.opt, func(i int) string { return cs.f(t, i) })
			require.Equal(t, cs.expected, result)
		})
	}
}

func TestFMap_ChainOfTransforms(t *testing.T) {
	t.Parallel()

	optInt := opt.Of(42)
	require.NotNil(t, optInt.Ptr())
	require.Equal(t, 42, *optInt.Ptr())

	optString := opt.FMap(optInt, strconv.Itoa)
	require.NotNil(t, optString.Ptr())
	require.Equal(t, "42", *optString.Ptr())

	optLength := opt.FMap(optString, func(s string) int { return len(s) })
	require.NotNil(t, optLength.Ptr())
	require.Equal(t, 2, *optLength.Ptr())
}

func TestFMapEx_SingleTransform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		opt      opt.Opt[int]
		f        func(t *testing.T, i int) (string, error)
		expected opt.Opt[string]
		err      error
	}{
		{
			"value 42",
			opt.Of(42),
			func(t *testing.T, i int) (string, error) {
				t.Helper()
				require.Equal(t, 42, i)

				return fmt.Sprintf("Value: %d", i), nil
			},
			opt.Of("Value: 42"),
			nil,
		},
		{
			"empty",
			opt.Empty[int](),
			func(t *testing.T, _ int) (string, error) { t.Helper(); t.Fail(); return "", nil },
			opt.Empty[string](),
			nil,
		},
		{
			"error",
			opt.Of(42),
			func(t *testing.T, i int) (string, error) {
				t.Helper()
				require.Equal(t, 42, i)

				return "", errTest
			},
			opt.Empty[string](),
			errTest,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result, err := opt.FMapEx(cs.opt, func(i int) (string, error) { return cs.f(t, i) })
			require.ErrorIs(t, err, cs.err)
			require.Equal(t, cs.expected, result)
		})
	}
}

var (
	errTestNoColon     = errors.New("no colon")
	errTestWrongLength = errors.New("wrong length")
)

func testSplitByColon(s string) ([]string, error) {
	if !strings.Contains(s, ":") {
		return nil, errTestNoColon
	}

	return strings.Split(s, ":"), nil
}

func testGetSecond(arr []string) (string, error) {
	if len(arr) != 2 {
		return "", errTestWrongLength
	}

	return arr[1], nil
}

func TestFMapEx_ChainOfTransforms(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		err      error
		phase    int
		opt      opt.Opt[string]
		expected opt.Opt[int]
	}{
		{"valid", nil, 0,
			opt.Of("value:42"), opt.Of(42)},

		{"valid with spaces", nil, 0,
			opt.Of("  value  :  42   "), opt.Of(42)},

		{"missing colon", errTestNoColon, 1,
			opt.Of("value42"), opt.Empty[int]()},

		{"wrong length", errTestWrongLength, 2,
			opt.Of("value:42:extra"), opt.Empty[int]()},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			optSplit, err := opt.FMapEx(cs.opt, testSplitByColon)
			if cs.phase == 1 {
				require.ErrorIs(t, err, cs.err)
				require.Nil(t, optSplit.Ptr())

				return
			}

			require.NoError(t, err)
			require.NotNil(t, optSplit.Ptr())

			optIndex, err := opt.FMapEx(optSplit, testGetSecond)
			if cs.phase == 2 {
				require.ErrorIs(t, err, cs.err)
				require.Nil(t, optIndex.Ptr())

				return
			}

			require.NoError(t, err)
			require.NotNil(t, optIndex.Ptr())

			optTrim, err := opt.FMapEx(optIndex,
				func(s string) (string, error) {
					return strings.TrimSpace(s), nil
				})
			require.NoError(t, err)
			require.NotNil(t, optTrim.Ptr())

			optInt, err := opt.FMapEx(optTrim, strconv.Atoi)
			require.NoError(t, err)
			require.NotNil(t, optInt.Ptr())
			require.Equal(t, cs.expected, optInt)
		})
	}
}
