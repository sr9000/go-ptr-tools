# No-Ptr: The Pointer

It's the simple helper module that allows you to generate a pointer to any given value.
The function `ptr.Of(value)` is useful when you need to pass a primitive value as a pointer, like `ptr.Of(42)`.

The function `ptr.Nil[T]()` is a drop-in replacement for `(*T)(nil)`.
This function is a first function of a package that allow you to avoid naming the pointer type explicitly.

> The whole idea of "No-Ptr" is to avoid using pointers explicitly in your codebase.

## Usage

[ptr_test.go](ptr_test.go)
