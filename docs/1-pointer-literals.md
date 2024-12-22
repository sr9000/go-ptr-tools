# Ptr-Tools: Pointer Literals

> Fully working code example provided at [pointer_to_literal_test.go](../examples/pointer_to_literal_test.go).

The first case is a limitation of Golang that does not allow getting the address of literal values. This often happens when filling a struct or passing argument to a function.

```go

// struct with pointers
type Bounds struct {
    Lower, Upper *float64
}

// func with pointers
func SplitString(s string, limit *int) []string { ... }
```

This limitation of Golang forcing to create a dummy variable. The function `ptr.New` allows to create a pointer in&#8209;place:

```go
bounds := Bounds{
    Lower: ptr.New(42.0), // provides &42.0
    Upper: nil,
}

// passing &3
SplitString("more than three words", ptr.New(3))
```



One more case that is not so common is when it's needed to wrap result of a function call with a pointer.

```go
func WordsLimit() int { ... }

// impossible in Golang
SplitString("more than three words", &WordsLimit())
```

This also can be fixed with `ptr.New`:

```go
SplitString("more than three words", ptr.New(WordsLimit()))
```

