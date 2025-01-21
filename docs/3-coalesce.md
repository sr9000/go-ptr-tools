# Ptr-Tools: Coalesce

> Fully working code example provided at [coalesce_test.go](../examples/coalesce_test.go).


Some values can be retrieved from several sources thus it is useful to choose one value using priorities. Here the
`Coalesce` function can help.

```go
func sourceA() *int { ... }
func sourceB() *int { ... }

func main() {
    res := ref.Coalesce(sourceA(), sourceB())
}
```

Also it is possible to provide default value so function `Else` should be used.

```go
func main() {
    fallback := 5
    res := ptr.Else(ref.Guaranteed(&fallback), sourceA(), sourceB())
}
```

## Recipe for slice of getters

> Fully working code example provided at [getters_coalesce_test.go](../examples/getters_coalesce_test.go).


Imagine there are several function returning a pointer, there is multiple ways to collect them. One of the possible
implementation provided in the following snippet of code.

```go
func CoalesceGetters[T any](getters ...func () *T) *T {
    var wg sync.WaitGroup
    
    wg.Add(len(getters))
    results := make([]*T, len(getters))
    
    for i, g := range getters {
        go func () {
            defer wg.Done()
            
            results[i] = g()
        }()
    }
    
    wg.Wait()

    return ptr.Coalesce(results...)
}
```

This implementation preserve order of getters and runs them in parallel.
