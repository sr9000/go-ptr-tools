package examples_test

import (
	"fmt"
	"strings"

	"github.com/sr9000/go-ptr-tools/ptr"
)

type Bounds struct {
	Lower, Upper *float64
}

func (l Bounds) String() string {
	var lower, upper any

	if l.Lower != nil {
		lower = *l.Lower
	}

	if l.Upper != nil {
		upper = *l.Upper
	}

	return fmt.Sprintf("{%v %v}", lower, upper)
}

func WordsLimit() int { return 3 }

func SplitString(s string, limit *int) []string {
	if limit == nil {
		return strings.Split(s, " ")
	}

	res := strings.SplitN(s, " ", *limit+1)
	if len(res) > *limit {
		res = res[:len(res)-1] // drop the last unsplit part
	}

	return res
}

func Example_fillingStruct() {
	limits := Bounds{
		Lower: ptr.New(42.0), // in-place pointer to literal
		Upper: nil,
	}

	fmt.Println(limits)
	// Output: {42 <nil>}
}

func Example_passingArg_wrappingResult() {
	text := "the quick brown fox jumps over the lazy dog"

	// passing pointer to literal as argument
	most5 := SplitString(text, ptr.New(5))

	// passing pointer to function result as argument
	mostLimit := SplitString(text, ptr.New(WordsLimit()))

	fmt.Println(len(most5), most5)
	fmt.Println(len(mostLimit), mostLimit)
	// Output:
	// 5 [the quick brown fox jumps]
	// 3 [the quick brown]
}
