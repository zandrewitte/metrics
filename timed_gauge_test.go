package metrics

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimedGaugeError(t *testing.T) {
	expectPanic(t, "NewTimedGauge_nil_callback", func() {
		NewTimedGauge("NewTimedGauge_nil_callback", nil)
	})
	expectPanic(t, "GetOrCreateTimedGauge_nil_callback", func() {
		GetOrCreateTimedGauge("GetOrCreateTimedGauge_nil_callback", nil)
	})
}

func TestTimedGaugeSerial(t *testing.T) {
	name := "GaugeSerial"
	n := 1.23
	var nLock sync.Mutex
	g := NewTimedGauge(name, func() (int64, float64) {
		nLock.Lock()
		defer nLock.Unlock()
		n++
		return 1234567890, n
	})
	for i := 0; i < 10; i++ {
		if _, nn := g.Get(); nn != n {
			t.Fatalf("unexpected gauge value; got %v; want %v", nn, n)
		}
	}

	// Verify marshalTo
	testMarshalTo(t, g, "foobar", "foobar 12.23 1234567890\n")

	// Verify big numbers marshaling
	n = 1234567899
	testMarshalTo(t, g, "prefix", "prefix 1234567900 1234567890\n")
}

func TestTimedGaugeConcurrent(t *testing.T) {
	name := "GaugeConcurrent"
	var n int
	var nLock sync.Mutex
	g := NewTimedGauge(name, func() (int64, float64) {
		nLock.Lock()
		defer nLock.Unlock()
		n++
		return time.Now().UnixMilli(), float64(n)
	})
	err := testConcurrent(func() error {
		_, nPrev := g.Get()
		for i := 0; i < 10; i++ {
			if _, n := g.Get(); n <= nPrev {
				return fmt.Errorf("gauge value must be greater than %v; got %v", nPrev, n)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
