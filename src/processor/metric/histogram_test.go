package metric

import (
	. "harvesterd/intf"
	"math"
	"runtime"
	"sync"
)

import . "launchpad.net/gocheck"

type HistogramSuite struct{}

var _ = Suite(&HistogramSuite{})

func (s *HistogramSuite) TestProcessInt(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	metric := NewHistogram("foo")

	var wait sync.WaitGroup
	var add = func() {
		for i := 1; i <= 10000; i++ {
			metric.Process(Record{"foo": i})
		}

		wait.Done()
	}

	count := 5
	for i := 0; i < count; i++ {
		go add()
	}

	wait.Add(count)
	wait.Wait()

	result := metric.GetValue().(map[string]interface{})
	c.Check(result["count"], Equals, int64(50000))
	c.Check(result["min"], Equals, 1.0)
	c.Check(result["max"], Equals, 10000.0)
	c.Check(result["mean"], Equals, 5000.5)
	c.Check(result["sum"], Equals, 2.50025e+08)
	c.Check(int(result["stddev"].(float64)), Equals, 2886)
}

func (s *HistogramSuite) TestProcessFloat64(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	metric := NewHistogram("foo")

	var wait sync.WaitGroup
	var add = func() {
		for i := 1; i <= 10000; i++ {
			metric.Process(Record{"foo": float64(i) / 1000.0})
		}

		wait.Done()
	}

	count := 5
	for i := 0; i < count; i++ {
		go add()
	}

	wait.Add(count)
	wait.Wait()

	result := metric.GetValue().(map[string]interface{})
	c.Check(result["count"], Equals, int64(50000))
	c.Check(result["min"], Equals, 0.001)
	c.Check(result["max"], Equals, 10.0)
	c.Check(result["mean"], Equals, 5.0004999812)
	c.Check(result["sum"], Equals, 250024.99906)
	c.Check(math.Floor(result["stddev"].(float64)*1000.0), Equals, math.Floor(2.886780196409398*1000.0))
}
