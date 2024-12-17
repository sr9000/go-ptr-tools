# No-Ptr: The Reference

This package brings the ability to pass a reference - an always valid pointer to a value.
The type `ref.Ref[T]` must be used as a function argument when you want to avoid unnecessary copying of the value.
But in contrast to the pointer, the reference is always valid and must not be checked for `nil`.

Summarizing the above, the type `ref.Ref[T]` plays a role of a contract for both the caller and the callee.

## Advantages

- The reference is always safe to use as a function argument.
- The reference preserves the pointer (`ref.OfPtr(ptr)`), so you can modify the original pointer value through the
  reference.

## Usage

[ref_test.go](ref_test.go)

## Caveats

Reference is not valid by default, so you must initialize your structs explicitly.

[ref_caveats_test.go](ref_caveats_test.go)
