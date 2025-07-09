# 1. Introduction

## Motivation

Go provides a minimal yet expressive core language with a strong emphasis on simplicity, readability, and explicitness. However, in day-to-day development, especially when dealing with pointers, value validation, and optional handling, idiomatic Go code can become overly verbose and repetitive.

Take the following common scenarios:

- Obtaining a pointer to a literal value
- Extracting a value from a result that includes an `ok` bool or `error`
- Coalescing potentially nil pointers to a first available one
- Ensuring a pointer is valid before dereferencing
- Passing a mutable value without allowing accidental `nil` access
- Providing default fallbacks for optional data

These operations often take multiple lines and copy-pasted idioms, reducing code readability and increasing room for subtle bugs.

This repository proposes a library for simplifying these patterns by introducing a minimal set of expressive, reusable utilities and types. It enables Go developers to work with pointers, references, and optional values using clear one-liners that preserve type safety and performance but improve code ergonomics.

The library is inspired by ideas from languages like Rust (e.g., `Option<T>`, references), Kotlin (e.g., safe call, elvis operator), and functional programming (monads). But it’s written **in idiomatic Go**, using standard language features available from Go 1.18+ (i.e., generics) without runtime cost.

It is not a framework and requires no integration or adoption ceremony. You can start using selected helpers in any Go codebase immediately — especially in the parts that deal with configuration, parsing, conversions, optional values, or pointer-heavy logic.

## Core Concepts: Pointer, Ref, Opt

The library introduces and works around three core data handling forms:

### Go Native Pointers: `*T`

Go has native support for pointers via `*T`. They are frequently used for:

- Representing optional values (e.g., `*int`)
- Allowing function arguments to be mutated
- Avoiding large-value copying (e.g. structs or slices)

But working with `*T` has drawbacks:
- You must check for `nil` everywhere before dereferencing
- Pointers to literals like `42` or `"hello"` cannot be written directly (need temp vars)
- There's no built-in mechanism to “validate” or “guarantee” non-nil pointers

This library provides helpers that address those without introducing runtime overhead — such as `ptr.New`, `ptr.FromOk`, `ptr.Coalesce`, and more.

### Always-Valid Pointer Wrapper: `Ref[T]`

`Ref[T]` is a generic struct that wraps a pointer and signaling its non-nil nature.

```go
type Ref[T any] struct {
  ptr *T
}
```

A `Ref[T]` is validated during construction (see `ref.New`) to ensure it always holds a valid pointer. It can be used:
- Safely, without repetitive `nil` checks
- As a lightweight replacement for passing `*T` when `nil` is not allowed
- To preserve mutability and avoid value copying

Advantages of `Ref[T]` include:
- Promised dereference safety
- Reduced repetitive checks
- Stronger API contracts (e.g., a function accepts `Ref[T]` to ensure a pointer is never nil)

However, `Ref[T]` comes with certain caveats. It should **not be embedded in structs** or used as a **zero-value** unless constructed explicitly. It is designed for **use in function arguments and return values only**, due to its construction-time guarantee.

### Lightweight Optional Value: `Opt[T]`

While Go does not have a built-in `Option` type, `Opt[T]` fills this role:

```go
type Opt[T any] struct {
  val T
  ok  bool
}
```

It represents an optional value without requiring pointers:
- `ok == true` means a value is present
- `ok == false` means the value is empty

Why use `Opt` instead of `*T`?
- No heap allocations — works entirely on stack
- Avoids pointer-based nil dereference bugs
- Helps when working with call signatures returning tuples (e.g., `(T, bool)` or `(T, error)`)

