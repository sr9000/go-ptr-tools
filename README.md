# go-ptr-tools

> Pointer is a new optional.

There are a lot of implementations of optional types in Golang like https://github.com/moznion/go-optional
or https://github.com/markphelps/optional.

But actually there is no lack in optional but in reference type. The `go-ptr-tools` package introduce a new type
`ref.Ref`.
Also, there is a set of widely used functions to work with pointers.

Additionally, a bunch of recipes to work with optional values using pointers are provided.
And some antipatterns are explained.

## How to work with repo

This repo is a library and a cookbook at the same time. Nothing to build or compile.

Bunch of Makefile commands are provided to run tests, benchmarks, and lint the code.

**Makefile Commands:**

- `lint`: Runs `golangci-lint` to lint the code and automatically fix issues.
- `test`: Runs all tests in the project using `go test`.
- `bench`: Runs benchmarks in the project using `go test -bench`.
- `clean`: Cleans the test cache using `go clean -testcache`.
- `all`: Cleans the test cache, lints the code, and runs all tests and benchmarks.

## Documentation and examples

1. [Pointer literals](docs/1-pointer-literals.md): Explains how to create pointers to literal values using the `ptr.New`
   function, which is useful for struct initialization and function arguments.
2. [Ref type](docs/2-reference-type.md): Describes the `ref.Ref` type, which ensures valid pointers are passed to
   functions, eliminating the need for nil checks.
3. [Coalesce optionals (pointers)](docs/3-coalesce.md): Demonstrates the `ref.Coalesce` function, which selects a value
   from multiple sources based on priority, and the `ptr.Else` function for providing default values.
4. [Optionals recipes](docs/4-optional-recipes.md): Provides recipes for working with optional values using pointers,
   including construction, validation, retrieval, and monad operations.
5. [Ref&Opt antipatterns](docs/5-anti-patterns.md): Lists common antipatterns when working with pointers and references
   in Go, along with solutions to avoid these pitfalls.
6. [*Extra: *T vs PT](docs/extra-t-vs-pt.md): Explains the difference between using `*T` directly and using `PT *T`in
   generics, highlighting the flexibility and strictness of each approach. Explains why library uses `*T` approach (
   spoiler — it's less restrictive).
