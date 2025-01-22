# Ptr-Tools: `*T` vs `PT *T` meaning and difference

> Fully working code example provided at [tpt_test.go](../examples/tpt_test.go).

⚠️ _This document created using LLM!_ ⚠️

Provided code example demonstrates the difference between directly using `*T` in a generic function and using `PT` as a
second type parameter constrained to be `*T`. Let's break it down step by step, explaining the behavior for each case
and highlighting the difference between `*T` and `PT`.

---

### **Key Concepts of `*T` and `PT`**

1. **`*T`**:
    - Directly represents a pointer to the generic type `T`.
    - Simple and convenient for working with pointers to types without requiring extra constraints or abstractions.

2. **`PT` (Pointer Type Parameter)**:
    - Introduced as a second type parameter (`PT *T`) in the generic function.
    - Allows extra abstraction, but it requires that `PT` explicitly match `*T`.
    - This can make it more restrictive and may require explicit type conversions when dealing with type aliases or
      structs.

---

### **Code Breakdown**

#### **Type Definitions**

```go
type Aint = *int // Aint is a type alias for *int
type Pint *int   // Pint is a new defined type, based on *int
type Sint struct { *int } // Sint is a struct with an embedded *int
```

- **`Aint`**: A type alias for `*int`. This means `Aint` is interchangeable with `*int`, and the Go compiler treats them
  as the same type.
- **`Pint`**: A distinct type that is based on `*int`. Even though `Pint` has the same underlying representation as
  `*int`, it is a different type and requires explicit conversion to/from `*int`.
- **`Sint`**: A struct that embeds a `*int`. Accessing the embedded pointer requires using the `int` field.

---

#### **Generic Functions**

##### **`foo`**

```go
func foo[T any](msg string, x *T) {
    fmt.Printf("foo %s: %T\n", msg, x)
}
```

- `foo` accepts a parameter `x` of type `*T`, where `T` is a generic type.
- It works with any type `T` and simply requires `x` to be a pointer (`*T`).
- No additional constraints are imposed, and it is flexible with type aliases, structs, or direct pointers.

##### **`bar`**

```go
func bar[T any, PT *T](msg string, x PT) {
    fmt.Printf("bar %s: %T\n", msg, x)
}
```

- `bar` introduces a second type parameter `PT` that must explicitly be a pointer to `T`.
- This makes `bar` more restrictive:
    - `PT` must exactly match `*T`.
    - If you use a distinct type like `Pint`, explicit conversion to `*T` is required.
    - This is less flexible than `foo`, as the Go compiler enforces stricter type matching.

---

#### **Example Function**

##### **1. Using a Pointer to an Integer**

```go
x := 42
foo("pointer", &x)
bar("pointer", &x)
```

- `&x` is a pointer to the integer `x`.
- For `foo`:
    - `T` is inferred as `int`, so `*T` becomes `*int`.
    - No issues arise because `*int` directly matches `*T`.
- For `bar`:
    - `T` is inferred as `int`, and `PT` is inferred as `*int`.
    - `&x` matches the expected `PT` (`*int`) without any explicit conversion.

Output:

```
foo pointer: *int
bar pointer: *int
```

---

##### **2. Using the `Pint` Wrapper**

```go
px := Pint(&x)
foo("wrapper", px)
bar("wrapper", (*int)(px)) // requires explicit conversion
```

- `Pint` is a distinct type based on `*int`.
- For `foo`:
    - `T` is inferred as `int`, and `*T` matches `Pint` because `Pint` has the same underlying representation as `*int`.
    - No explicit conversion is required.
- For `bar`:
    - `PT` must be exactly `*T`, but `Pint` is a distinct type, so it doesn't directly match `*T`.
    - You must explicitly convert `Pint` to `*int` using `(*int)(px)`.

Output:

```
foo wrapper: *int
bar wrapper: *int
```

---

##### **3. Using the `Aint` Alias**

```go
ax := Aint(px)
foo("alias", ax)
bar("alias", ax)
```

- `Aint` is a type alias for `*int`.
- For both `foo` and `bar`:
    - Since `Aint` is treated as identical to `*int` by the compiler, no explicit conversion is needed.
    - Both `foo` and `bar` work seamlessly with `ax`.

Output:

```
foo alias: *int
bar alias: *int
```

---

##### **4. Using the `Sint` Struct**

```go
sx := Sint{&x}
foo("struct", sx.int) // requires field access
bar("struct", sx.int) // requires field access
```

- `Sint` is a struct with an embedded `*int` field.
- To pass the pointer stored in `Sint` to `foo` or `bar`, you must explicitly access the `int` field (`sx.int`).
- For both `foo` and `bar`:
    - `T` is inferred as `int`, and the `*int` stored in `sx.int` matches `*T`.
    - No explicit conversion is needed, but you need field access.

Output:

```
foo struct: *int
bar struct: *int
```

---

### **Key Differences Between `*T` and `PT`**

| Feature                   | `foo` (`*T` directly)                           | `bar` (`PT *T`)                                    |
|---------------------------|-------------------------------------------------|----------------------------------------------------|
| **Type Matching**         | Works directly with any `*T`.                   | `PT` must explicitly match `*T`.                   |
| **Type Aliases**          | Works seamlessly with type aliases like `Aint`. | Works seamlessly with type aliases like `Aint`.    |
| **Distinct Types**        | Works with distinct types like `Pint`.          | Requires explicit conversion for distinct types.   |
| **Structs with Pointers** | Requires explicit field access for structs.     | Requires explicit field access for structs.        |
| **Flexibility**           | More flexible, simpler to use.                  | More restrictive, requires stricter type matching. |

---

### **Summary**

1. **`foo` (with `*T`)**:
    - Flexible and works directly with any type of pointer, including type aliases (`Aint`), distinct types (`Pint`),
      and embedded pointers (`Sint`).
    - No explicit conversions are required except for accessing struct fields.

2. **`bar` (with `PT *T`)**:
    - More restrictive and enforces that `PT` matches `*T` exactly.
    - Requires explicit type conversion for distinct types like `Pint`, even though they are based on `*int`.
    - Works seamlessly with type aliases like `Aint`, as they are treated as identical to `*int`.

The primary difference lies in the **strictness** of type matching. `foo` is simpler and more flexible, while `bar`
enforces stricter constraints due to the additional type parameter `PT`.

So this package uses `*T` to keep code simple and flexible, allowing for easier integration with various types and
functions.
