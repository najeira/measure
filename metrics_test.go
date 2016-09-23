package measure

import (
	"log"
	"strconv"
	"testing"
	"time"
)

var enablePrint = false

func printf(format string, a ...interface{}) {
	if enablePrint {
		log.Printf(format, a...)
	}
}

func TestMeasure(t *testing.T) {
	const key = "test"

	Reset()
	m := Start(key)
	m.Stop()

	stats := GetStats()
	if stats == nil || len(stats) < 1 {
		t.Fatal("GetStats returns nil")
	}

	if stats[0].Count != 1 {
		t.Errorf("Stats.Count got %d expect %d", stats[0].Count, 1)
	}
}

func TestMeasureMulti(t *testing.T) {
	const key = "test_multi"
	const loop = 100

	Reset()
	for i := 0; i < loop; i++ {
		m := Start(key)
		time.Sleep(time.Microsecond)
		m.Stop()
	}

	stats := GetStats()
	if stats == nil || len(stats) < 1 {
		t.Fatal("GetStats returns nil")
	}

	if stats[0].Count != loop {
		t.Errorf("Stats.Count got %d expect %d", stats[0].Count, loop)
	}

	if stats[0].Min > stats[0].Max {
		t.Errorf("Stats.Min %f > Stats.Max %f", stats[0].Min, stats[0].Max)
	}
	if stats[0].Min > stats[0].Avg {
		t.Errorf("Stats.Min %f > Stats.Max %f", stats[0].Min, stats[0].Avg)
	}
	if stats[0].Max < stats[0].Avg {
		t.Errorf("Stats.Min %f < Stats.Max %f", stats[0].Max, stats[0].Avg)
	}
	if stats[0].Sum < stats[0].Max {
		t.Errorf("Stats.Sum %f < Stats.Max %f", stats[0].Sum, stats[0].Max)
	}
}

func TestMeasureSort(t *testing.T) {
	const loop = 100

	Reset()
	for i := 0; i < loop; i++ {
		m := Start(strconv.Itoa(i))
		time.Sleep(time.Microsecond)
		m.Stop()
	}

	stats := GetStats()
	if stats == nil || len(stats) < 1 {
		t.Fatal("GetStats returns nil")
	}

	stats.SortAsc("sum")
	for i := 0; i < len(stats)-1; i++ {
		n, m := stats[i], stats[i+1]
		if n.Sum > m.Sum {
			t.Fatal("SortAsc fail")
		}
		printf("%f <= %f", n.Sum, m.Sum)
	}

	stats.SortDesc("sum")
	for i := 0; i < len(stats)-1; i++ {
		n, m := stats[i], stats[i+1]
		if n.Sum < m.Sum {
			t.Fatal("SortDesc fail")
		}
		printf("%f >= %f", n.Sum, m.Sum)
	}
}

func TestMeasureMetrics(t *testing.T) {
	m1 := NewMetrics()
	m2 := NewMetrics()

	s1 := m1.Start("test")
	s2 := m2.Start("test")

	s1.Stop()
	s2.Stop()

	if len(m1.metrics) != 1 {
		t.Error("invalid len")
	} else if _, ok := m1.metrics["test"]; !ok {
		t.Error("invalid map")
	}

	if len(m2.metrics) != 1 {
		t.Error("invalid len")
	} else if _, ok := m2.metrics["test"]; !ok {
		t.Error("invalid map")
	}

	log.Print(m1.GetStats())
	log.Print(m2.GetStats())
}

func BenchmarkMeasure(b *testing.B) {
	const key = "test"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m := Start(key)
		m.Stop()
	}
}

func BenchmarkMeasureDisabled(b *testing.B) {
	const key = "test"
	Disabled = true
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m := Start(key)
		m.Stop()
	}
	Disabled = false
}
