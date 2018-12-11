package main

import (
	"bytes"
	"encoding/binary"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"hash/fnv"
)

var (
	hash   = fnv.New64a()
	strBuf bytes.Buffer // Used for hashing.
	intBuf = make([]byte, 8)
)

// hashNameAndLabels returns a hash value of the provided name string and all
// the label names and values in the provided labels map.
//
// Not safe for concurrent use! (Uses a shared buffer and hasher to save on
// allocations.)
func hashNameAndLabels(name string, labels prometheus.Labels) uint64 {
	hash.Reset()
	strBuf.Reset()
	strBuf.WriteString(name)
	hash.Write(strBuf.Bytes())
	binary.BigEndian.PutUint64(intBuf, model.LabelsToSignature(labels))
	hash.Write(intBuf)
	return hash.Sum64()
}

type Exporter struct {
	Counters   *CounterContainer
	Gauges     *GaugeContainer
	Histograms *HistogramContainer
}

func NewExporter() *Exporter {
	return &Exporter{
		Counters:   NewCounterContainer(),
		Gauges:     NewGaugeContainer(),
		Histograms: NewHistogramContainer(),
	}
}

type CounterContainer struct {
	Elements map[uint64]prometheus.Counter
}

func NewCounterContainer() *CounterContainer {
	return &CounterContainer{
		Elements: make(map[uint64]prometheus.Counter),
	}
}

func (c *CounterContainer) Get(metricName string, labels prometheus.Labels, help string) (prometheus.Counter, error) {
	hash := hashNameAndLabels(metricName, labels)
	counter, ok := c.Elements[hash]
	if !ok {
		counter = prometheus.NewCounter(prometheus.CounterOpts{
			Name:        metricName,
			Help:        help,
			ConstLabels: labels,
		})
		if err := prometheus.Register(counter); err != nil {
			return nil, err
		}
		c.Elements[hash] = counter
	}
	return counter, nil
}

type GaugeContainer struct {
	Elements map[uint64]prometheus.Gauge
}

func NewGaugeContainer() *GaugeContainer {
	return &GaugeContainer{
		Elements: make(map[uint64]prometheus.Gauge),
	}
}

func (c *GaugeContainer) Get(metricName string, labels prometheus.Labels, help string) (prometheus.Gauge, error) {
	hash := hashNameAndLabels(metricName, labels)
	gauge, ok := c.Elements[hash]
	if !ok {
		gauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        metricName,
			Help:        help,
			ConstLabels: labels,
		})
		if err := prometheus.Register(gauge); err != nil {
			return nil, err
		}
		c.Elements[hash] = gauge
	}
	return gauge, nil
}

type HistogramContainer struct {
	Elements map[uint64]prometheus.Histogram
}

func NewHistogramContainer() *HistogramContainer {
	return &HistogramContainer{
		Elements: make(map[uint64]prometheus.Histogram),
	}
}

func (c *HistogramContainer) Get(metricName string, labels prometheus.Labels, help string) (prometheus.Histogram, error) {
	hash := hashNameAndLabels(metricName, labels)
	histogram, ok := c.Elements[hash]

	if !ok {
		histogram = prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:        metricName,
			Help:        help,
			ConstLabels: labels,
			Buckets: []float64{0.125, 0.25, 0.5, 1, 2, 5, 7, 10, 15, 20, 25, 30, 40, 50, 60, 70,
				80, 90, 100, 200, 300, 400, 500, 1000,
				2000, 5000, 10000, 30000, 60000},
		})
		if err := prometheus.Register(histogram); err != nil {
			return nil, err
		}
		c.Elements[hash] = histogram
	}
	return histogram, nil
}
