# 3. Working with Pointers

Go pointers (`*T`) are a powerful and efficient way to convey optional values, avoid value copying, and enable mutability. However, the language deliberately offers minimal syntactic sugar for pointer creation and usage. This often results in boilerplate code, bloat of `if` conditions, or confusing workarounds when dealing with optional or dynamic values.

This section introduces utility functions from the `ptr` package that simplify everyday pointer operations while staying within Go idioms. These helpers reduce verbosity, clarify intent, and eliminate error-prone patterns from pointer-heavy code.

## Creating Pointers from Values

In Go, taking the address of a literal is not always straightforward:

```go
x := 42
ptr := &x // extra variable, can't do &42 directly
```

Library provides `ptr.Of` — a concise one-liner for such use cases:

```go
p := ptr.Of(42)           // *int
s := ptr.Of("hello, Go")  // *string
```

**Definition:**

```go
func Of[T any](v T) *T {
  return &v
}
```

Creates a pointer to a value `T`. This is ideal for:

- Passing pointer arguments derived from constants or literals,
- Inline use in composite literals or configuration structs,
- Creating well-typed test data.

## Safe Pointer Creation with `ptr` Helpers

### Conditional Value Acquisition

When working with conditionally valid values in Go—such as those returned with an `ok` flag, an `error`, or compared against their zero value—creating pointers often involves verbose if-statements. Here are some common examples:

```go
val, ok := m["key"]
var ptr1 *string
if ok {
  ptr1 = &val
}

val, err := os.ReadFile("file.txt")
var ptr2 *[]byte
if err == nil {
  ptr2 = &val
}
```

The `ptr` package provides  `ptr.FromOk` and `ptr.FromErr` to simplify these common patterns:

```go
val, ok := m["key"]
ptr1 := ptr.FromOk(val, ok)

ptr2 := ptr.FromErr(os.ReadFile("file.txt"))
```

- `ptr.FromOk(val, ok)` returns a pointer to `val` only if `ok` is `true`
- `ptr.FromErr(val, err)` returns a pointer to `val` only if `err == nil`

### Checking for zero values

Another common scenario is assigning a pointer only if a value is **not equal** to its zero value. For example:

```go
var namePtr *string

if name = GetNameFromEnv(); name != "" {
  namePtr = &name
}
```

Instead of writing conditional logic, use `ptr.FromZero`:

```go
namePtr := ptr.FromZero(GetNameFromEnv())
```

The `ptr.FromZero` function returns a pointer to the value if it’s not equal to its zero value:

```go
ptr.FromZero("")      // nil
ptr.FromZero("admin") // *string
ptr.FromZero(0)       // nil
ptr.FromZero(42)      // *int
```

## Pointer Coalescing and Fallbacks

When dealing with multiple sources that might yield `nil` pointers (such as config overrides or fallback data), it is common to need the **first non-nil pointer** available.

Instead of verbose condition chains:

```go
if env.Proxy != nil {
  return env.Proxy
} else if cfg.Proxy != nil {
  return cfg.Proxy
}
// ...
return defaultProxy
```

Use `ptr.Coalesce` that finds and returns the **first non-nil pointer** in the provided list. It returns nil if all provided pointers are nil. This makes fallback chains clear and expressive.

```go
effectiveProxy := ptr.Coalesce(env.Proxy, cfg.Proxy, sys.Proxy)
```

