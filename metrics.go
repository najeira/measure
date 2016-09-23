package measure

import (
	"sort"
	"strings"
	"sync"
	"time"

	mt "github.com/rcrowley/go-metrics"
)

const (
	Key   = "key"
	Count = "count"
	Sum   = "sum"
	Min   = "min"
	Max   = "max"
	Avg   = "avg"
	Rate  = "rate"
	P95   = "p95"
)

var (
	Disabled bool

	defaultMetrics *Metrics
)

func init() {
	defaultMetrics = NewMetrics()
}

type Measure struct {
	key     string
	start   time.Time
	metrics *Metrics
}

func Start(key string) Measure {
	return defaultMetrics.Start(key)
}

func (m Measure) Stop() {
	if Disabled {
		return
	}
	m.metrics.Update(m.key, m.start)
}

type Metrics struct {
	mu      sync.RWMutex
	metrics map[string]mt.Timer
}

func GetStats() StatsSlice {
	return defaultMetrics.GetStats()
}

func Reset() {
	defaultMetrics.Reset()
}

func NewMetrics() *Metrics {
	return &Metrics{metrics: make(map[string]mt.Timer)}
}

func (m *Metrics) Start(key string) Measure {
	if Disabled {
		return Measure{}
	}
	return Measure{key: key, start: time.Now(), metrics: m}
}

func (m *Metrics) Update(key string, start time.Time) {
	m.mu.RLock()
	t, ok := m.metrics[key]
	m.mu.RUnlock()

	if !ok {
		m.mu.Lock()
		t, ok = m.metrics[key]
		if !ok {
			t = mt.NewTimer()
			m.metrics[key] = t
		}
		m.mu.Unlock()
	}

	t.Update(time.Since(start))
}

func (m *Metrics) GetStats() StatsSlice {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(StatsSlice, 0, len(m.metrics))
	for key, t := range m.metrics {
		stats := Stats{
			Key:   key,
			Count: t.Count(),
			Sum:   float64(t.Sum()) / float64(time.Millisecond),
			Min:   float64(t.Min()) / float64(time.Millisecond),
			Max:   float64(t.Max()) / float64(time.Millisecond),
			Avg:   t.Mean() / float64(time.Millisecond),
			Rate:  t.Rate1(),
			P95:   t.Percentile(0.95) / float64(time.Millisecond),
		}
		result = append(result, stats)
	}
	return result
}

func (m *Metrics) Reset() {
	m.mu.Lock()
	m.metrics = make(map[string]mt.Timer)
	m.mu.Unlock()
}

type Stats struct {
	Key   string  `csv:"key" json:"key"`
	Count int64   `csv:"count" json:"count"`
	Sum   float64 `csv:"sum" json:"sum"`
	Min   float64 `csv:"min" json:"min"`
	Max   float64 `csv:"max" json:"max"`
	Avg   float64 `csv:"avg" json:"avg"`
	Rate  float64 `csv:"rate" json:"rate"`
	P95   float64 `csv:"p95" json:"p95"`
}

type StatsSlice []Stats

func (s StatsSlice) SortAsc(key string) {
	s.sortBy(key, true)
}

func (s StatsSlice) SortDesc(key string) {
	s.sortBy(key, false)
}

func (s StatsSlice) sortBy(key string, asc bool) {
	p := StatsSorter{stats: s, key: key}
	if asc {
		sort.Sort(p)
	} else {
		sort.Sort(sort.Reverse(p))
	}
}

type StatsSorter struct {
	stats []Stats
	key   string
}

func (p StatsSorter) Len() int {
	return len(p.stats)
}

func (p StatsSorter) Less(i, j int) bool {
	n, m := p.stats[i], p.stats[j]
	switch p.key {
	case Key:
		return strings.Compare(n.Key, m.Key) < 0
	case Count:
		return n.Count < m.Count
	case Min:
		return n.Min < m.Min
	case Max:
		return n.Max < m.Max
	case Avg:
		return n.Avg < m.Avg
	case Rate:
		return n.Rate < m.Rate
	case P95:
		return n.P95 < m.P95
	}
	return n.Sum < m.Sum
}

func (p StatsSorter) Swap(i, j int) {
	p.stats[i], p.stats[j] = p.stats[j], p.stats[i]
}
