# Ptr-Tools: Anti Patterns

> Fully working code example provided at [anti_patterns_test.go](../examples/anti_patterns_test.go).

## Missing Pointer Checks

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

## `ref.Ref` Declarations

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

## `ref.Ref` Struct Fields

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
    n, err := ref.New(s.N)
    if err != nil {
        return err
    }

    // use n and other fields
}

func main() {
	err := Consumer(S{})  // safe to pass zero value
}
```
