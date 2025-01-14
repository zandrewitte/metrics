package metrics_test

import (
	"bytes"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
)

func ExampleSet() {
	// Create a set with a counter
	s := metrics.NewSet()
	sc := s.NewCounter("set_counter")
	sc.Inc()
	s.NewGauge(`set_gauge{foo="bar"}`, func() float64 { return 42 })
	s.NewTimedGauge(`set_timed_gauge{foo="bar"}`, func() (int64, float64) { return 1234567890, 42 })

	// Dump metrics from s.
	var bb bytes.Buffer
	s.WritePrometheus(&bb)
	fmt.Printf("set metrics:\n%s\n", bb.String())

	// Output:
	// set metrics:
	// set_counter 1
	// set_gauge{foo="bar"} 42
	// set_timed_gauge{foo="bar"} 42 1234567890
}
