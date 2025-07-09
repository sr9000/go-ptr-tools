# 4. Always-Present References

> [!IMPORTANT]
>
> `Ref[T]` is **not** a replacement for `*T` in all situations — instead, it’s a tool for structuring your code around the assumption of *non-optional shared data*.

In Go, most of the time we rely on `*T` pointers to avoid copying. But when you're *certain* a pointer must be non-nil, unchecked use of `*T` can still lead to errors: for example the pointer can be zero-initialized or unset by mistake.

`Ref[T]` is a thin wrapper over `*T` that guarantees the pointer is always valid and cannot be nil after construction. Use `Ref[T]` when:

1. You want to **share access** to a nontrivial data structure — e.g., config, model, buffer, etc. — across multiple layers or components, and nils are not expected or allowed.

2. You want to **mutate a shared value** before/after passing it:
   ```go
   func AdjustTimeout(conf ref.Ref[Config]) {
     conf.Ptr().Timeout += 3 * time.Second
   }
   ```

3. The underlying type **should not be copied**:
   - Large struct (including slices, maps, nested fields)
   - Stateful type (e.g., with sync primitives)
   - Frequently reused object (to avoid GC churn)

4. You want to guarantee that something **has been initialized**, even if its value came from optional sources.

## Ref Constructors

To enforce the invariant that `Ref[T]` always wraps a valid pointer, construction must go through one of the provided constructors:

### FromPtr

```go
func FromPtr[T any](ptr *T) (Ref[T], error)
```
This constructor performs a runtime nil check and returns an error if the pointer is nil. Safe for handling dynamic values or uncertain sources.

**Example:**

```go
r, err := ref.FromPtr(ptr)
if err != nil {
  return err // or fallback
}
```

### Guaranteed

```go
func Guaranteed[T any](ptr *T) Ref[T]
```

A **"trust me"** constructor. You assert that `ptr` is non-nil and the function skips validation.

Use only when you are absolutely certain:

```go
r := ref.Guaranteed(&myVal)
```

### Of

> [!CAUTION]
>
> Beware you'll get a reference to a **COPY** of a value. If you need a reference to an original value use `ref.Guaranteed(&originalValue)` as described above.

```go
func Of[T any](v T) Ref[T]
```

Creates a `Ref[T]` by taking a pointer to a literal or value. Handy when you're working with value literals:

```go
ref := ref.Of(42)
```

You get a valid, permanent reference without needing to create a variable first.

## Accessing Data in Ref

Once a `Ref[T]` is created, it's safe to access its value without checking for nil. Two key methods are provided:

1. `Val()`
2. `Ptr()`

```go
func (r Ref[T]) Val() T
```

Returns a **copy** of the value stored in the reference.

```go
func (r Ref[T]) Ptr() *T
```

Returns the **underlying pointer** so it can be passed to functions accepting `*T`.

## Working with Ref in APIs

Use `Ref[T]` in APIs when:

- You **require the input** to be a valid, non-nil value (e.g., validation done upstream)
- You want to avoid passing copies of large structs or nested data
- You want to encapsulate "non-optional, mutable" semantics

### Good Example: Function Arguments

```go
type Config struct {
    Timeout time.Duration
    Buffer  []byte
    Limits  []int
    // many more fields...
}
```

You're passing it into several layers of your system, and some layers update it:

```go
func NormalizeLimits(cfg ref.Ref[Config]) {
    if len(cfg.Ptr().Limits) == 0 {
        cfg.Ptr().Limits = []int{100, 500, 1000}
    }
}
```

This communicates at the type level that the input *must be present*.

### Anti-Pattern: `Ref[T]` as Struct Field

Avoid:

```go
type Person struct {
	Name ref.Ref[string]  // BAD pattern!
}
```

Reasons:
- Cannot enforce construction safety
- Zero-value `Ref[T]{}` is invalid (contains nil pointer)
- Can lead to runtime panics if accessed without appropriate constructor

## Interoperability with Pointers

`Ref[T]` integrates seamlessly with functions and APIs that deal with raw pointers:

- If a caller provides a `*T`, you can convert to `Ref[T]` safely with `ref.FromPtr` or `ref.Guaranteed`.
- To return `*T` for interop, use `r.Ptr()`.

### Use With Fallback Logic

You can combine pointer utilities with `ref.Else` to ensure a guaranteed reference when multiple optional pointers are available:

```go
f := ptr.Else(ref.Of("fallback"), inputA, inputB)
```

Returns a `Ref[string]` from the first non-nil pointer, or uses the fallback value.

### Summary of Conversion Patterns

| Input Type | To `Ref[T]`           | Notes                      |
| ---------- | --------------------- | -------------------------- |
| `*T`       | `ref.FromPtr(ptr)`    | Safe, returns error if nil |
| `*T`       | `ref.Guaranteed(ptr)` | Unsafe, assumes non-nil    |
| `T`        | `ref.Of(v)`           | Creates owned reference    |
| `Ref[T]`   | `r.Ptr()`             | For pointer-based interop  |
| `Ref[T]`   | `r.Val()`             | For direct value access    |

`Ref[T]` is not a replacement for pointers — it’s a structured way of ensuring non-nil semantics *when nil is not an option*.

