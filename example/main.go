// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"time"

	"github.com/uber-go/tally"
)

type printStatsReporter struct{}

// NewCactusStatsReporter returns a new StatsReporter
func newPrintStatsReporter() tally.StatsReporter {
	return &printStatsReporter{}
}

func (r *printStatsReporter) ReportCounter(name string, tags map[string]string, value int64) {
	fmt.Printf("%s %d\n", name, value)
}

func (r *printStatsReporter) ReportGauge(name string, tags map[string]string, value int64) {
	fmt.Printf("%s %d\n", name, value)
}

func (r *printStatsReporter) ReportTimer(name string, tags map[string]string, interval time.Duration) {
	fmt.Printf("%s %d\n", name, interval)
}

func main() {
	reporter := newPrintStatsReporter()
	rootScope := tally.NewScope("", nil, reporter, time.Second)
	subScope := rootScope.SubScope("requests")

	bighand := time.NewTicker(time.Millisecond * 2300)
	littlehand := time.NewTicker(time.Millisecond * 10)

	measureThing := rootScope.Gauge("thing")
	tickCounter := subScope.Counter("ticks")

	// Spin forever, watch report get called
	go func() {
		for {
			select {
			case <-bighand.C:
				measureThing.Update(42)

			case <-littlehand.C:
				tickCounter.Inc(1)
			}
		}
	}()

	select {}
}