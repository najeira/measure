# measure

## Usage

### Measure

Add measure to your code.

```go
import "github.com/najeira/measure"

func foo() {
    defer measure.Start("foo").Stop()

    // your code

}
```

or

```go
...
m := measure.Start("foo")
// your code
m.Stop()
...
```

### Stats

Get statistics.

```go
stats := measure.GetStats()
stats.SortDesc("sum")
for _, stat := range stats {
    fmt.Printf("%s = %f\n", stat.Key, stat.Sum)
}
```

Reset statistics.

```go
measure.Reset()
```

## License

MIT
