# 6. Functional Style: Monadic Optional Handling

Go encourages explicit and imperative handling of absence: declaring `*T` and checking for `nil` before proceeding. While safe and idiomatic, this can lead to deeply nested code, broken pipelines, and repetitive patterns — especially where value transformation is expected.

This section introduces a **functional style abstraction** using monadic operations on pointers (`*T`). These abstractions enable safe composition and transformation of optional values (*i.e., values wrapped in `*T` that may be `nil`*) using standardized, one-liner-friendly helpers.


## What is a Monad?

In the context of pointers, a *monad* is a pattern that allows for expressing computations over a value that may or may not exist (`*T`), **without having to manually check for `nil`** at every step.

Instead of writing this:

```go
if input != nil {
  out := f(*input)
  // more logic...
}
```

Or even worse — move presence checking into calling function:

```go
func f(input *Type) Result {
  if input == nil {
    // absence branch
  } else {
    // presence branch
  }
}
```

You can write calling function as usual and let all the checks for the monad:

```go
out := ptr.Apply(input, f)
```

This monadic interface lets you isolate presence-checking in a reusable utility and focus on the transformation logic itself.


## `Apply` Functions

The `Apply` family of functions encapsulates a common idiom:
- **If** a value (`*T`) is present
  - **Then** apply a function (`T → R`)
  - **And Return** the result wrapped in `*R`

- **Else** return `nil`

### Basic Form

```go
func Apply[R1, T1 any](
  t1 *T1,
  fn func(t1 T1) R1,
) *R1
```

**Example:**

```go
upper := ptr.Apply(userName, strings.ToUpper)
```

**Explanation:**

- If `t1 == nil` → return `nil`
- If `t1 != nil` → invoke `fn(*t1)`, return pointer to result

This pattern avoids:
- Manual `nil` checks
- Variable shadowing for output
- Deep nesting in data pipelines

### Available Forms Of `Apply`

The library includes numerous generated `Apply` variants. These allow:

- Handling multiple pointer inputs
- Mapping to multiple return values
- Supporting `void` (side-effect only) functions
- Wrapping functions with `error` returns
- Inserting `context.Context` where needed

They are discussed in details in section *"Naming Convention for Operation Signatures"*. 

#### Examples:

| Wrapper Name    | Wraps                               |
| --------------- | ----------------------------------- |
| `Apply2`        | `func(t1, t2) R`                    |
| `ApplyVoid`     | `func(t1)`                          |
| `Apply3CtxErr`  | `func(ctx, t1, t2, t3) (R, error)`  |
| `Apply11`       | `func(t1) (R1, R2)`                 |
| `Apply95CtxErr` | `func(ctx, t1…t9) → (r1…r5, error)` |

All variants behave the same:
- If **any pointer input is `nil`**, the function is **not called**,
- The result(s) are returned as pointers where appropriate.

## `Monad` Wrappers

While `Apply` calls a function immediately, the `Monad` family wraps a function and returns a new function that takes pointer arguments and returns pointer results.

In other words:

```go
wrappedF := Monad(f)
var r *R = wrappedF(ptr)
```

This is ideal for:
- **Reusing pointer-aware logic**
- **Chaining** transformations
- **Passing function pipelines** to higher-order functions

### Basic Form

```go
func Monad[R1, T1 any](
  fn func(T1) R1,
) func(*T1) *R1
```

**Example #1:**

The function that takes an optional string, returns an optional uppercase version.

```go
toUpper := ptr.Monad(strings.ToUpper)

var res []*string
for _, u := range users {
  u.FirstName = toUpper(u.FirstName)   // safely applied
  u.LastName = toUpper(u.LastName)     // even if value are missing
  u.MiddleName = toUpper(u.MiddleName) // without repetative if x != nil
}
```

**Example #2:**

Composing monads into a transformation pipeline.

```go
trim := ptr.Monad(strings.TrimSpace)
lower := ptr.Monad(strings.ToLower)

normalize := func(s *string) *string {
  return lower(trim(s))
}
```

You can apply `normalize` to any optional `*string` — safely and predictably.

## Naming Convention for Operation Signatures

All `Apply` and `Monad` functions follow a structured naming pattern to describe their behavior, matching a consistent, predictable grammar.

Pattern is `[Apply|Monad]<N>[<M>|Void][Ctx][Err]`, where:

| Part   | Meaning                                                      |
| ------ | ------------------------------------------------------------ |
| `N`    | Number of arguments; omitted if N = 1 and M < 2              |
| `M`    | Number of output results (>1); omitted if M = 1              |
| `Void` | There is no return values (used for side-effects); used when M = 0 |
| `Ctx`  | First parameter is a `context.Context`                       |
| `Err`  | Last return value is an `error`                              |

**Examples:**

| Function Name      | Interpreted As                                               |
| ------------------ | ------------------------------------------------------------ |
| `Apply`            | 1 input, 1 output: `func(t1 T1) R1`                          |
| `Monad2`           | 2 inputs, 1 output: `func(t1 T1, t2 T2) R1`                  |
| `Monad4VoidCtxErr` | 4 inputs, no return, uses `context.Context`, returns `error` |
| `Apply12`          | 1 input, 2 outputs: `func(t1 T1) (R1, R2)`                   |
| `Apply95CtxErr`    | 9 inputs, 5 outputs, accepts `context.Context`, returns `error`.<br />(Full signature provided as a separate code snippet below) |
```go
// `fn` argument for Apply95CtxErr
func(ctx context.Context,
  t1 *T1, ..., t9 *T9,
) (
  r1 *R1, ..., r5 *R5,
  err error,
)
```


This suffix system helps you quickly identify the correct helper for your use case — whether simple one-value mapping or high-arity, context-aware operation with error handling. All are implemented with Go generics, and return `nil` for all outputs if any pointer in the input set is nil.
