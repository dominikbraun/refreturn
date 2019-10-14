# refreturn
> Find functions that return a reference and cause allocations.

When a function allocates a value and returns a reference to it, the value has to escape to the heap. This is slower and puts pressure on the garbage collector.

**This may be optimized: You can avoid the heap allocation and allow the function to be inlined.**

---

### Example: a simple constructor

```go
struct Coffee {
    Type string
}

func New() *Coffee {
    c := Coffee{
        Type: "espresso"
    }
    return &c
}
```

0. [Download refreturn](https://github.com/dominikbraun/refreturn/releases) and copy the binary into your project's root for example.
1. Run `./refreturn <directory>` and you'll see that `New` returns a reference.
2. Check if the returned value is being created in the function.
3. This is true for our `c` variable.
4. Optimize the function like so:

```go
func New() *Coffee {
    var c Coffee
    return new(&c)
}

func new(c *Coffee) *Coffee {
    c.Type = "espresso"
    return c
}
```

`New()` is now merely a wrapper which allocates the instance. The "real work" will be done in `new()`.

**This will allow mid-stack inlining as described in [this blog post](https://blog.filippo.io/efficient-go-apis-with-the-inliner/).**
