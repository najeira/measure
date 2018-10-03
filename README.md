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

// print stats in CSV format
for _, s := range stats {
	fmt.Fprintf(w, "%s,%d,%f,%f,%f,%f,%f,%f\n",
		s.Key, s.Count, s.Sum, s.Min, s.Max, s.Avg, s.Rate, s.P95)
}
```

Reset statistics.

```go
measure.Reset()
```

### Metrics

You can handle multiple metrics.

```go
var metricsA = measure.NewMetrics()
var metricsB = measure.NewMetrics()

func foo() {
    defer metricsA.Start("foo").Stop()
}

func bar() {
    defer metricsB.Start("bar").Stop()
}
```

## License

MIT
