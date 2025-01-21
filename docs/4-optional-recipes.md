# Ptr-Tools: Optional Recipes

> Fully working code example provided at [optional_recipes_test.go](../examples/optional_recipes_test.go).

## Construct Optional

**1. Optional.Of() value from a variable**

```go
var n int
opt = &n
```

**2. Optional.Of() literal value**

```go
opt = ptr.New(42)
```

**3. Optional.Empty()**

```go
var opt *int = nil // in case of variable
return nil         // in case of function
```

## Validate Optional

**4. Optional.IsPresent()**

```go
if opt != nil { ... }
```

**5. Optional.IsEmpty()**

```go
if opt == nil { ... }
```

## Retrieve Value

**6. Optional.Get()**

```go
if opt != nil {
    val = *opt
}
```

**7. Optional.OrElse()**

```go
ref = ptr.Else(ref.Literal(42), opt)
```

## Monad Operations

**8. Optional.Coalesce()**

```go
res := ptr.Coalesce(optA, optB)
```

**9. Optional.Map()**

```go
func Map[T any, R any](opt *T, f func (T) R) *R {
    if opt == nil {
        return nil
    }
    
    return ptr.New(f(*opt))
}
```

Or inside a function:

```go
var r *R
if opt != nil {
    r = ptr.New(f(*opt))
}
```
