package metrics

import (
	"fmt"
	"io"
)

// NewTimedGauge registers and returns gauge with the given name, which calls f
// to obtain the gauge value and timestamp associated with the value.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// f must be safe for concurrent calls.
//
// The returned gauge is safe to use from concurrent goroutines.
//
// See also FloatCounter for working with floating-point values.
func NewTimedGauge(name string, f func() (int64, float64)) *TimedGauge {
	return defaultSet.NewTimedGauge(name, f)
}

// TimedGauge is a float64 gauge.
//
// See also Counter, which could be used as a gauge with Set and Dec calls.
type TimedGauge struct {
	f func() (int64, float64)
}

// Get returns the current value for g.
func (g *TimedGauge) Get() (int64, float64) {
	return g.f()
}

func (g *TimedGauge) marshalTo(prefix string, w io.Writer) {
	t, v := g.f()
	if float64(int64(v)) == v {
		// Marshal integer values without scientific notation
		fmt.Fprintf(w, "%s %d %d\n", prefix, int64(v), t)
	} else {
		fmt.Fprintf(w, "%s %g %d\n", prefix, v, t)
	}
}

// GetOrCreateTimedGauge returns registered gauge with the given name
// or creates new gauge if the registry doesn't contain gauge with
// the given name.
//
// name must be valid Prometheus-compatible metric with possible labels.
// For instance,
//
//   - foo
//   - foo{bar="baz"}
//   - foo{bar="baz",aaa="b"}
//
// The returned gauge is safe to use from concurrent goroutines.
//
// Performance tip: prefer NewGauge instead of GetOrCreateGauge.
//
// See also FloatCounter for working with floating-point values.
func GetOrCreateTimedGauge(name string, f func() (int64, float64)) *TimedGauge {
	return defaultSet.GetOrCreateTimedGauge(name, f)
}
