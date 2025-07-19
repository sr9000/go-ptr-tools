# 8. Value-Based Optional Type

> [!TIP]
>
> Do **not** use `Opt[T]`.

`Opt[T]` is a generalized value container that represents presence or absence of a value in a pure, allocation-free form.

```go
type Opt[T any] struct {
  val T
  ok  bool
}
```

This is effectively just a `(T, bool)` tuple with helper methods attached. No magic, no allocation, and no special contract beyond what `if ok` already expresses. The type exists **only** to simplify certain awkward Go patterns — particularly where `(T, bool)` is produced but not immediately consumed.

## Reasoning about `Opt[T]`

Use `Opt[T]` **only if all of the following apply**:

- You need to represent an optional value
- You want to **avoid heap allocations** (i.e. you’d like to avoid `*T`)
- The `(T, bool)` pair is **produced in one place** and **consumed later**
- You don't want to propagate separate `val` and `ok` variables everywhere

In Go, there’s no native way to store or return a `(T, bool)` pair as a real value. You’re forced to either:
- Reunpack the values repeatedly
- Or use a custom struct

While `Opt[T]` solves this, it is **not** meant to replace pointers, it’s a niche tool designed to carry state cleanly — and nothing more.

## Working with `Opt[T]`

### Constructing

**Literal Value:**

```go
opt := opt.Of(42) // value is present
```

**Empty Value (default zero):**

```go
var opt opt.Opt[int] // value is .ok == false
```

**Conditional Value:**

```go
opt := opt.FromOk(v, ok) // manually wrap
```

### From Other Patterns

```go
opt.FromErr(v, err)    // from (T, error)
opt.FromZero(v)        // if v != zero(T)
opt.FromPtr(ptr)       // from *T
```

Each version uses clear, consistent behavior — present only when valid.

### Checking for Value

```go
if opt.IsMissing() {
  return SomeDefaultValue // exit early
}
```

Or:

```go
if v, ok := opt.Get(); ok {
	use(v) // guarded call
}
```

### Filtering Chains

Instead of spreading presence checks everywhere, use `Coalesce` or `Else`:

```go
firstAvailable := opt.Coalesce(optA, optB, optC)
```

Returns the first present value, or empty if all are empty.

```go
v := opt.Else(42, optA, optB, optC)
```

Returns the first present value, or a final fallback (42 for this example).

## Monad Support

`Opt[T]` supports the same monadic operation patterns as `*T`, so you can write cleaner pipelines without branching or unpacking.

```go
opt.Apply(optA, func(a T) R) // returns Opt[R]
```

Applies `fn` to the inner value — only if present.

```go
normalize := opt.Monad(func(s string) string {
  return strings.ToLower(s)
})
```

Defines `normilize` wrapper that can be used later like:

```go
out := normalize(optA)
```

Equivalent to:

```go
out := opt.Apply(optA, strings.ToLower)
```

## When `Opt[T]` Adds Real Value

> TL;DR: **don’t reach for `Opt[T]` by default.** Reach for it when standard Go makes you reach too far.

### Async Code with Heap Allocation Awareness

**Scenario:**

You start multiple goroutines that each return an optional result. You **want to collect these into a single slice** without allocating a `*T` in every goroutine.

Using pointers in goroutines often causes heap allocations — each goroutine forces the captured locals to escape.

**Step #1: fetch some required data**

```go
data := make([]opt.Opt[Data], len(keys))
g, ctx := errgroup.WithContext(ctx)

for i, k := range keys {
  i, k := i, k // capture loop vars
  g.Go(func() error {
    val, ok, err := fetchData(ctx, k)
    if err != nil {
      return err
    }
    
    data[i] = opt.FromOk(val, ok) // wrap into optional
    return nil
  })
}

// Wait for all fetches to complete
if err := g.Wait(); err != nil {
  return err
}
```

**Step #2: get actual items**

```go
items := make([]Item, len(keys))
g, ctx := errgroup.WithContext(ctx)

for i := range keys {
  i := i // capture loop variable
  g.Go(func() error {
    var (
      item Item
      err  error
    )

    if data, ok := results[i].Get(); ok {
      item, err = getItem(ctx, keys[i], data)
    } else {
      item, err = legacyItem(ctx, keys[i])
    }

    if err != nil {
      return err
    }

    items[i] = item
    return nil
  })
}

if err := g.Wait(); err != nil {
  return nil, err
}
```

**Step #3: collect final map**

```go
// Build final map
m := make(map[Key]Item, len(keys))
for i, item := range items {
  m[keys[i]] = item
}

```

Here, using `Opt[T]` keeps present/missing status together **in-place**, without any heap allocation per result or awkward use of a twin slice of bools. Using `*T` in this context would generally require heap allocation due to pointer capture in goroutines.

### Multi-Argument Optional Input Where Values Might Escape

**Scenario:**

You want to pass multiple conditionally available inputs into a pipeline:

```go
optA := opt.FromOk(aVal, aOk)
optB := opt.FromOk(bVal, bOk)
optC := opt.FromOk(cVal, cOk)

res := process(optA, optB, optC)
```

Here, using `Opt[T]` keeps present/missing status together **in-place**, without bloating function arguments with "ok" booleans like that:

```go
func process(a Val, aOk bool, b Val, bOk bool, c Val, cOk bool) Result {
  if aOk && bOk && cOk {
    res := doSomeStuff(aVal, bVal, cVal)
    // ...
  }
}
```

With `Opt[T]` arguments function defined like that:

```go
func process(a, b, c opt.Opt[Val]) Result {
  if res, ok := opt.Apply3(a, b, c, doSomeStuff).Get(); ok {
    // ...
  }
}
```

> [!CAUTION]
>
> Using `*Val` in this context may help but requires to be **EXTRA** accurate to avoid heap allocation. At the same time `Opt[T]` always stays on stack due to its value-based nature and do not require **manual** extra validation using escape analysis.

## Recommendations

> [!TIP]
>
> Do not use `Opt[T]`.

While `Opt[T]` provides a clean way to store `(T, ok)` pairs without heap allocation, most use cases can be expressed **more simply, more idiomatically, and with better clarity** using native Go patterns: `*T`, `T, bool`, or `T, error`.

This section walks through common scenarios where you might **reach for `Opt[T]`**, explains why it’s **unnecessary or even counterproductive**, and offers **idiomatic alternatives**.

> [!TIP]
>
> Filter out early and collect only valid values.

**Temptation:**

You want to collect values from a source where some keys might be missing. It feels natural to wrap each result in `Opt[T]`, store in a slice, and later use only ones which are present:

```go
var filtered []opt.Opt[string]
for _, key := range keys {
  filtered = append(filtered, opt.FromOk(fetch(key)))
}

// Process only present values
for _, o := range filtered {
  if v, ok := o.Get(); ok {
    fmt.Println("Using:", v)
  }
}
```

**Why `Opt[T]` is overkill:**

The presence info is available immediately — you don’t need to preserve `(val, ok)` for later. Cleaner, fewer types, zero mental overhead.

```go
var values []string
for _, key := range keys {
  if v, ok := fetch(key); ok {
    values = append(values, v)
  }
}

// Use values directly
for _, v := range values {
  fmt.Println("Using:", v)
}
```

> [!TIP]
>
> Use a `*T` if you have no special concerns about heap allocations.
>

**Temptation:**

You're writing a function that might or might not act on a value if provided. So you select `Opt[T]`:

```go
func IDontCareAboutHeapAllocations(name opt.Opt[string]) {
  if name.IsPresent() {
    // Run with name
  } else {
    // Fallback or ignore
  }
}
```

**Why `Opt[T]` is overkill:**

Using non-idiomatic constructions must be explained. Here the heap allocations is not a concern. So native approach is more idiomatic and do not require the caller to import `Opt[T]`.

```go
func IDontCareAboutHeapAllocations(name *string) {
  if name != nil {
    // Run with name
  } else {
    // Fallback or ignore
  }
}
```

> [!TIP]
>
> Or use an `ok bool` to avoid heap allocation.

```go
func WhyNotToCareAboutHeapAllocations(name string, ok bool) {
  if ok {
    // Run with name
  } else {
    // Fallback or ignore
  }
}
```

> [!TIP]
>
> Or even allow the caller to decide.

```go
func ZeroHeapAllocations(name string) {
  // Run with name
}
```

Caller side:

```go
if name, ok := getName(); ok {
  ZeroHeapAllocations(name)
} else {
  // Fallback or ignore
}
```

One top-level decision about presence — and zero new types.

> [!TIP]
>
> Replacing missing values with a default one.

**Temptation:**

You’d like to replace missing/malformed data with a default later, so something like:

```go
var fields []opt.Opt[string]
for _, id := range ids {
  v, ok := lookup[id]
  fields = append(fields, opt.FromOk(v, ok))
}

for _, f := range fields {
  fmt.Println("Field:", opt.Else("unknown", f))
}
```

**Why it's overkill:**

You’re storing fully constructed values. You might as well store final values from the beginning.

```go
var fields []string
for _, id := range ids {
  if v, ok := lookup[id]; ok {
    fields = append(fields, v)
  } else {
    fields = append(fields, "unknown")
  }
}

for _, f := range fields {
  fmt.Println("Field:", f)
}
```

No boxing/unboxing. It’s already the right type.
