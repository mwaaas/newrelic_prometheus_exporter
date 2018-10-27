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
}

func NewExporter() *Exporter {
	return &Exporter{
		Counters:   NewCounterContainer(),
		Gauges:     NewGaugeContainer(),
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