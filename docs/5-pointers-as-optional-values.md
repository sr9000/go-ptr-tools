# 5. Pointers as Optional Values

The idiomatic way to represent optional values is using `*T` pointers. When `nil`, the value is considered absent; else the value is considered present.

- **Create Optional from Variable:** take an address.

  ```go
  opt := &value
  ```

- **Create Optional from Literal:** use `ptr.Of(...)` for a one-liner shortcut. Useful in function arguments, map values, or optional configurations.

  ```go
  opt := ptr.Of(42)         // *int
  name := ptr.Of("Alice")   // *string
  ```

- **Represent Absence:** simply use `nil`.

  ```go
  var absent *int = nil
  ```

- **Check Presence:** compare against `nil`.

  ```go
  if opt != nil {
    fmt.Println(*opt)
  }
  ```

- **Extract Value (if needed):** dereference a pointer or use `*opt` directly in tightly scoped function parameters. Do not forget to guard dereference by presence check.

  ```go
  val := 0
  if opt != nil {
    val = *opt
  }
  ```

- **Coalescing (First Non-Nil):** use `ptr.Coalesce(...)` to find the first available value.

  ```go
  result := ptr.Coalesce(optA, optB, optC)
  ```

- **Fallback with Literal (Non-Nil Guarantee):** use `ptr.Else(...)` combined with `ref.Literal(...)` to safely convert an uncertain chain into a guaranteed `Ref[T]`. This allows you to treat optional inputs declaratively, while guaranteeing that the result is always safe to use.

  ```go
  f := ptr.Else(
      ref.Literal("default"),
      optA, optB, optC,
  )
  ```

## Rest Utilities

| Function            | Purpose                                |
| ------------------- | -------------------------------------- |
| `ptr.FromZero(...)` | Converts zero-checkable values to `*T` |
| `ptr.FromOk(...)`   | Converts a value and a `bool` to `*T`  |
| `ptr.FromErr(...)`  | Converts a value and `error` to `*T`   |