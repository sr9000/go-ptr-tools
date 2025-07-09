# go-ptr-tools

[![Go Lint](https://github.com/sr9000/go-ptr-tools/actions/workflows/lint.yml/badge.svg)](https://github.com/sr9000/go-ptr-tools/actions/workflows/lint.yml)
[![Go Test](https://github.com/sr9000/go-ptr-tools/actions/workflows/test.yml/badge.svg)](https://github.com/sr9000/go-ptr-tools/actions/workflows/test.yml)
[![Go Coverage](https://github.com/sr9000/go-ptr-tools/wiki/coverage.svg)](https://raw.githack.com/wiki/sr9000/go-ptr-tools/coverage.html)

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

## Getting Started

This library is a zero-dependency utility package designed to be embedded directly into your Go projects. It contains simple building blocks and helper functions around safe value handling using pointers, references, and optional types, with no runtime configuration or special environment required.

There are **no binaries to build**, no frameworks to install, and no background processes involved. The repository provides a set of utility packages and a **Makefile** to run tests, benchmarks, and lint checks easily.

**Installation**:

To start using the utilities, simply import the relevant packages into your Go project.

```bash
go get https://github.com/sr9000/go-ptr-tools
```

**Get sources repo**:

```bash
git clone https://github.com/sr9000/go-ptr-tools.git
```

**Makefile Commands**:

A Makefile is included to facilitate common development tasks like linting, running tests, generating benchmarks, or clearing the test cache. This helps keep your workflow clean and consistent.

Here are the available **Makefile targets**:

| Command      | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| `make lint`  | Runs `golangci-lint` to lint the code and optionally auto-fix style issues. |
| `make test`  | Runs all Go tests in the project using `go test`.            |
| `make cover` | Runs tests and generates a test coverage report.             |
| `make bench` | Runs benchmarks using the `go test -bench` command.          |
| `make clean` | Clears the Go test cache using `go clean -testcache`.        |
| `make all`   | Performs `clean`, `lint`, `test`, and `bench` in sequence.   |



## Documentation

1. [Introduction](docs/1-introduction.md)  
   1.1 Motivation  
   1.2 Core Concepts: Pointer, Ref, Opt

2. [Getting Started](docs/2-getting-started.md)  
   2.1 Installation  
   2.2 Get sources repo  
   2.3 Makefile Commands

3. [Working with Pointers](docs/3-working-with-pointers.md)  
   3.1 Creating Pointers from Values  
   3.2 Safe Pointer Creation with `ptr` Helpers  
   3.3 Pointer Coalescing and Fallbacks

4. [Always-Present References](docs/4-always-present-references.md)  
   4.1 Ref Constructors  
   4.2 Accessing Data in Ref  
   4.3 Working with Ref in APIs  
   4.4 Interoperability with Pointers

5. [Pointers as Optional Values](docs/5-pointers-as-optional-values.md)  
   5.1 Create Optional from Variable  
   5.2 Create Optional from Literal  
   5.3 Represent Absence  
   5.4 Check Presence  
   5.5 Extract Value (if needed)  
   5.6 Coalescing (First Non-Nil)  
   5.7 Fallback with Literal (Non-Nil Guarantee)  
   5.8 Rest Utilities

6. [Functional Style: Monadic Optional Handling](docs/6-functional-style-monadic-optional-handling.md)  
   6.1 What is a Monad?  
   6.2 `Apply` Functions  
   6.3 `Monad` Wrappers  
   6.4 Naming Convention for Operation Signatures

7. [Common Antipatterns & Pitfalls](docs/7-common-antipatterns-and-pitfalls.md)  
   7.1 Using `bool` to represent optionality  
   7.2 Omitting nil checks when dereferencing `*T`  
   7.3 Declaring `Ref[T]` as a struct field  
   7.4 Wrapping `ref.Ref[T]` in `any` or creating `*Ref[T]` (pointer to a reference)  
   7.5 Returning `Ref[T]` in combination with `ok bool` or `error`  
   7.6 Reconstructing `Ref[T]` unnecessarily  
   7.7 Using Coalesce with fallback value  
   7.8 The desire to change where reference is pointing  
   7.9 Assuming `Ref[T]` is safe for concurrent access  
   7.10 Marshaling `Ref[T]

8. [Value-Based Optional Type](docs/8-value-based-optional-type.md)  
   8.1 Reasoning about `Opt[T]`  
   8.2 Working with `Opt[T]`  
   8.3 Monad Support  
   8.4 When `Opt[T]` Adds Real Value  
   8.5 Recommendations

9. [*Extra: *T vs PT](docs/extra-t-vs-pt.md): Explains the difference between using `*T` directly and using `PT *T`in
   generics, highlighting the flexibility and strictness of each approach. Explains why library uses `*T` approach (
   spoiler — it's less restrictive).
