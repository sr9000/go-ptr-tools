# Ptr-Tools: Anti Patterns

> Fully working code example provided at [anti_patterns_test.go](../examples/anti_patterns_test.go).

## 1. Using `bool` To Check Validity

**Problem**: `bool` used to check validity of a function result.

```go
func find(xs []int, x int) (int, bool) {
    for i, v := range xs {
        if v == x {
            return i, true
        }
    }
    
    return 0, false
}
```

The `bool` type is not a good way to check validity. It is possible to forget to check the result.

**Solution**: Use a pointer to indicate the absence of a result.

```go
func find(xs []int, x int) *int {
    for i, v := range xs {
        if v == x {
            return &i
        }
    }
    
    return nil
}
```

As an alternative, it is possible to use `error` type to indicate the invalid result.

```go
func find(xs []int, x int) (*int, error) {
    for i, v := range xs {
        if v == x {
            return &i, nil
        }
    }
    
    return nil, errors.New("not found")
}
```

## 2. Missing Pointer Checks

**Problem**: Missing `nil` pointer checks.

The argument is usually "this code guarantees that pointer is valid, so there is no need to check it".

```go
func plusOne(n *int) {
    *n++
}
```

It's a common mistake, because nothing can stop the caller from passing a `nil` pointer.
Also, nothing in a function signature indicates that the pointer is expected to be valid.

**Solution**: Use `ref.Ref` type to separate *possibly nil* pointers from *valid* pointers.

```go
func plusOneToRef(r ref.Ref[int]) {
    *r.Ptr()++
}
```

If it is not a case for a `ref.Ref` type, then it is a case for `nil` pointer checks.

```go
func plusOneToPtr(n *int) {
    if n != nil {
        *n++
    }
}
```

## 3. `ref.Ref` Declarations

**Problem**: `ref.Ref` used for a variable declaration.

```go
var n ref.Ref[int]
```

The `n` reference is not initialized and keeps `nil` pointer. This reference is invalid and dangerous.
Declaration of a references is against the concept of the `ref.Ref` type.

**Solution**: Use function to initialize reference variable.

Every reference must be initialized with `New`, `Guaranteed`, or `Literal` functions.
Or from any other function that returns a reference.

```go
n := ref.Literal(42)
```

## 4. `ref.Ref` Struct Fields

**Problem**: `ref.Ref` used for a struct field.

```go
type S struct {
    N ref.Ref[int]
}
```

It is possible to create a zero value of the struct. It follows the same problem as with the variable declaration.

**Solution**: Use pure values if possible.

Often, it is possible to use a value instead of a reference. This is a better approach.

```go
type S struct {
    N int
}
```

Sometimes, it is necessary to modify an original value, so there is no way to avoid `ref.Ref` in this case.
Thus, the mechanics to suppress manual struct initialization should be implemented.

For example, a constructor function with functional options. Yes, it is still possible to create an invalid struct,
<u>but</u> an existence of a constructor function indicates to the caller that the struct should not be created manually.

```go
type S struct {
    N        ref.Ref[int]
    isValidN bool
}

func (s *S) isValid() bool {
    return s.isValidN // or any other validation
}

func NewS(options ...func (ref.Ref[S])) (S, error) {
    var s S
    
    for _, opt := range options {
        opt(ref.Guaranteed(&s))
    }
    
    if !s.isValid() {
        return S{}, errors.New("invalid S")
    }
    
    return s, nil
}
```

If this approach looks like an overkill, then it is better to use a pointer. Yes, it requires double conversion
(ref to pointer and pointer to ref), but it keeps code safe. Even with the cost to spawn errors on `nil` pointers.

```go
type S struct {
    N *int
    // lots of other fields
}

func Consumer(s S) error {
    n, err := ref.New(s.N)  // convert pointer back to ref
    if err != nil {
        return err
    }

    // use n and other fields
}

func Producer1() {
	err := Consumer(S{})  // safe to pass zero value
}

func Producer2(n ref.Ref[int]) {
    err := Consumer(S{N: n.Ptr()})  // safe to pass pointer
}
```

## 5. Using Pointer To `ref.Ref` (Or Store `ref.Ref` As `any`)

**Problem**: `ref.Ref` stored as `any` value or pointer to `ref.Ref`.

```go
func foo(r *ref.Ref[int]) {
    *r.Ptr()++
}

func bar(r any) {
    if n, ok := r.(ref.Ref[int]); ok {
        *n.Ptr()++
    }
}
```

The `ref.Ref` type is a reference type and not required to be wrapped in another reference type.

**Solution**: Use `ref.Ref` directly.

```go
func foobar(r ref.Ref[int]) {
    *r.Ptr()++
}
```

## 6. Wrap `any` Value Or A Pointer With `ref.Ref`

**Problem**: `any` type or pointer used in `ref.Ref`.

```go
func foo(r ref.Ref[any]) {
    *r.Ptr() = 42
}

func bar(r ref.Ref[*int]) {
    **r.Ptr() = 42
}
```

The `ref.Ref` type is a reference type itself and should not be used to wrap other reference types.

**Solution**: Use `ref.Ref` directly.

```go
func foobar(r ref.Ref[int]) {
    *r.Ptr() = 42
}
```

Even when it's needed to wrap a field or variable to change it.

```go
type S struct {
	n int
}

func foo(n ref.Ref[int]) {
    *n.Ptr()++
}

func bar(s ref.Ref[S]) {
    foo(ref.Guaranteed(&s.Ptr().n))
}
```

## 7. Using `ref.Ref` For Optional Results

**Problem**: `ref.Ref` used to wrap optional results.

```go
func minimum(xs []*int) ref.Ref[*int] {
    if len(xs) == 0 {
        return ref.Literal(nil) // no minimum for empty list
    }

	// assume that all pointers are valid
    min := &xs[0]
    for i := range xs {
        if *xs[i] < **min {
            min = &xs[i]
        }
    }
	
    return ref.Guaranteed(min) // reference to slice element itself
}
```

The `ref.Ref` type represents a valid pointer. This type must be treated like a value (e.g. `int`).
Thus, it is not suitable for optional results.

**Solution**: Use pointers directly.

```go
func minimum(xs []*int) **int {
    if len(xs) == 0 {
        return nil
    }

    min := &xs[0]
    for i := range xs {
        if *xs[i] < **min {
            min = &xs[i]
        }
    }

    return min
}
```

## 8. Extra Actions To Make A Copy

**Problem**: Extra actions to make a copy of a `ref.Ref` value.

```go
func foobar() {
	var n, m int
	// ...
	r1 := ref.Guaranteed(&n)
	// ...
	r2 := ref.Guaranteed(r1.Ptr())  // extra actions to make a copy
	r1 = ref.Guaranteed(&m)         // change r1
}
```

The `ref.Ref` is safe to copy. It is basically a pointer inside a struct.

**Solution**: Use an assignment operator.

```go
func foobar() {
    var n, m int
    // ...
    r1 := ref.Guaranteed(&n)
    // ...
    r2 := r1                 // safe to copy
    r1 = ref.Guaranteed(&m)  // reassignment doesn't affect r2
}
```

## 9. Unprotected Concurrent Access

**Problem**: `ref.Ref` used in concurrent access.

```go
func foobar() {
    var n int
	var wg sync.WaitGroup
	
	wg.Add(2)
    r := ref.Guaranteed(&n)
    
    go func() {
		defer wg.Done()
        *r.Ptr()++
    }()
    
    go func() {
		defer wg.Done()
        *r.Ptr()++
    }()
	
	wg.Wait()
}
```

The `ref.Ref` is not safe to concurrent access by default.

**Solution**: Use synchronization primitives as usual.

```go
func foobar() {
    var n int
    r := ref.Guaranteed(&n)
    
    var mu sync.Mutex
	var wg sync.WaitGroup
	
	wg.Add(2)
    
    go func() {
		defer wg.Done()
        mu.Lock()
		defer mu.Unlock()
        *r.Ptr()++
    }()
    
    go func() {
		defer wg.Done()
        mu.Lock()
		defer mu.Unlock()
        *r.Ptr()++
    }()
	
	wg.Wait()
}
```

## 10. Marshaling `ref.Ref` Values

**Problem**: `ref.Ref` values are marshaled as a part of a struct.

```go
type S struct {
    N ref.Ref[int]
}

func saveJSON(s S) {
    data, _ := json.Marshal(s)
    // save data
}
```

The `ref.Ref` type is a struct and not suitable for marshaling.

**Solution**: Use a value directly.

```go
type S struct {
    N ref.Ref[int]
	// lots of other fields
}

func saveJSON(s S) {
    x := struct {
        N int
		// lots of other fields
    }{
        N: s.N.Val(),
		// lots of other fields
    }

    data, _ := json.Marshal(x)
    // save data
}
```
