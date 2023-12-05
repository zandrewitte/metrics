package metrics_test

import (
	"fmt"
	"runtime"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

func ExampleTimedGauge() {
	// Define a timed gauge exporting the number of goroutines.
	var g = metrics.NewTimedGauge(`goroutines_count`, func() (int64, float64) {
		return time.Now().UnixMilli(), float64(runtime.NumGoroutine())
	})

	// Obtain gauge value.
	fmt.Println(g.Get())
}

func ExampleTimedGauge_vec() {
	for i := 0; i < 3; i++ {
		// Dynamically construct metric name and pass it to GetOrCreateGauge.
		name := fmt.Sprintf(`timed_metric{label1=%q, label2="%d"}`, "value1", i)
		iLocal := i
		metrics.GetOrCreateTimedGauge(name, func() (int64, float64) {
			return 1234567890, float64(iLocal + 1)
		})
	}

	// Read counter values.
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf(`timed_metric{label1=%q, label2="%d"}`, "value1", i)
		tm, n := metrics.GetOrCreateTimedGauge(name, func() (int64, float64) { return 0, 0 }).Get()
		fmt.Printf("%f %d\n", n, tm)
	}

	// Output:
	// 1.000000 1234567890
	// 2.000000 1234567890
	// 3.000000 1234567890
}
