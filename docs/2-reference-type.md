# Ptr-Tools: Reference Function Arguments

> Fully working code example provided at [ref_args_test.go](../examples/ref_args_test.go).


Pointer is the only way to mutate a value passed into a function. Pointers must be checked on `nil` before every usage
and this is very annoying even only valid pointers are passed into a function.

```go
func plusOne(n *int) {
    if n != nil { // checks happend ...
        *n++
    }
}

func main() {
    var n int
    plusOne(&n) // every
    plusOne(&n) // signle
    plusOne(&n) // time
}
```

Of course, it is possible to write function without checks... Oh, snap!

> panic: runtime error: invalid memory address or nil pointer dereference

The type `ref.Ref` provide a way to separate *valid* pointers from *possibly nil* pointers. Regular pointers are
*possibly nil* pointers and plays kinda like optional type in other languages. The `ref.Ref` representing a *valid*
pointers.

Yep, there is no way to guarantee validity using Golang type system but when an argument has type `ref.Ref` it's
saying "Hi there, I'm holding a valid pointer! There is no need to check me against nil!". The callee is responsible for
providing only valid pointers using `ref.Ref` type.

If an address of an actual variable is passed then it's a case for `ref.Guaranteed` function:

```go
func plusOne(r ref.Ref[int]) {
    *r.Ptr()++ // no more checks
}

func main() {
    var n int
    plusOne(ref.Guaranteed(&n))
    plusOne(ref.Guaranteed(&n))
    plusOne(ref.Guaranteed(&n))
}
```

If an optional value must be passed as reference then it's a case for `ref.New` function:

```go
func tryPlusOne(n *int) error {
    r, err := ref.New(n)
    if err != nil {
        return err
    }

    plusOne(r)

    return nil
}
```
