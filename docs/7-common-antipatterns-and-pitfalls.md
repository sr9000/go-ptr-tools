
# 7. Common Antipatterns & Pitfalls

When working with pointers `*T` and validated references `Ref[T]`, it’s easy to fall into unsafe, inconvenient or unidiomatic usage patterns — especially when trying to simplify pointer logic or avoid boilerplate manually.

Below is a list of common mistakes and better alternatives.

#### 🚫 Using `bool` to represent optionality

```go
val, ok := parse(input)
if ok {
  return &val
}
return nil
```

While returning `(T, bool)` is common in idiomatic Go (e.g. `map[key]` lookup), it’s awkward when you want optional **references** or to propagate that state across function boundaries.

✅ Correct: use `ptr.FromOk` — one-liner that produces a `*T` only if present.

```go
val, ok := parse(input)
return ptr.FromOk(val, ok)
```
#### 🚫 Omitting nil checks when dereferencing `*T`

```go
func foobar(reference *Type) { // assumed reference cannot be nil
  fmt.Println(*ptr) // panic if ptr == nil
}
```

✅ Correct: use `Ref[T]` instead of bare pointer if the value should always be valid.

```go
func foobar(r ref.Ref[Type]) {
  fmt.Println(r.Val())
}
```

#### 🚫 Declaring `Ref[T]` as a struct field

```go
type Config struct {
  Port ref.Ref[int]  // ❌ Invalid if zero-value
}
```

This is dangerous because `Ref[T]` constructed without validation will panic on access.

✅ Correct: use bare pointers and extra `nil` checks in that particular case. `Ref[T]` allowed to use in function signatures only.

```go
type Config struct {
  Port *int
}

func foobar(cfg Config) {
  if cfg.Port != nil {
    // ...
  }
}
```

#### 🚫 Wrapping `ref.Ref[T]` in `any` or creating `*Ref[T]` (pointer to a reference)

```go
r := ref.Of("value")

var anything any = r // try to reference a Ref
refPtr := &r         // try to reference the same value
```

Don’t hide `Ref[T]` behind a generic, you lose compile-time type safety and make runtime behavior unpredictable. Also `Ref[T]` is already a reference value. Making a pointer to a `Ref[T]` is an unnecessary extra-step.

✅ Correct: use `Ref[T]` explicitly or a known interface with methods. If needs a copy — use an assignment operator.

```go
r := ref.Of("value")
refCopy := r // reference the same value

// or do not use Ref at all: replace with it known interface with methods
var s Stringer = ...
```

#### 🚫 Returning `Ref[T]` in combination with `ok bool` or `error`

```go
func getValue() (ref.Ref[int], error)
```

This defeats the entire purpose of `Ref[T]`: guaranteeing presence. If there’s an error or need for optionality, return a `*T` directly.

✅ Correct:
```go
func getValue() (*int, error)
```

#### 🚫 Reconstructing `Ref[T]` unnecessarily

```go
r2 := ref.Guaranteed(r1.Ptr())
```

✅ Correct: use an assignment. `Ref[T]` is a simple struct and safe to copy directly.
```go
r2 := r1
```

#### 🚫 Using Coalesce with fallback value

```go
effectivePortPtr := ptr.Coalesce(
  envPort,
  filePort,
  ptr.Of(8080),
)
```

✅ Correct: use ref, when an alternative is guaranteed.
```go
effectivePortRef := ptr.Else(ref.Of(8080), envPort, filePort)
```

Also it is good style to avoid extra wrappers around primitives, so even better approach:
```go
effectivePortVal := 8080

if envPort != nil {
  effectivePortVal = *envPort
} else if filePort != nil {
  effectivePortVal = *filePort
}
```

Still simple and idiomatic. Common sense here is "do not use extra tools if possible".

#### 🚫 The desire to change where reference is pointing

```go
r := ref.Of(42)
for {
  patch(&r)
}
```

✅ Correct: use `r.Ptr()` and change underlying value directly.

```go
r := ref.Of(42)
for {
  patch(r)
}

func patch(r ref.Ref[int]) {
    *r.Ptr() = ... // safe to dereference without checking for nil
}
```

If the address is matter, return new reference value:

```go
r := ref.Of(42)
for {
  r = patch(r)
}

func patch(r ref.Ref[int]) ref.Ref[int] {
    if r.Val() > 0 {
        return r
    } else {
        return extractSpecialReference()
    }
}
```

#### 🚫 Assuming `Ref[T]` is safe for concurrent access

```go
r := ref.Of(42)
for {
  go func() {
    // race condition
    if r.Val() < 0 {
      r = ...
    }
  }
}
```

`Ref[T]` is only a struct — not a thread-safe primitive.

✅ Correct: use sync primitives.
```go
mu.Lock()
if r.Val() < 0 {
  r = ...
}
mu.Unlock()
```

#### 🚫 Marshaling `Ref[T]`

```go
r := ref.Of(42)
json.Marshal(r)
```

The library does **not** provide built-in JSON or binary marshaling for `Ref[T]`. If doing that you've got an empty object `{}` because all fields of `Ref[T]` are not exported.

✅ Correct: extract pointer or value before serialization.
```go
r := ref.Of(42)
json.Marshal(r.Ptr()) // ok, marshaling called on a pointer to a value
json.Marshal(r.Val()) // also ok, a copy of a value are made beffore serialization
```
