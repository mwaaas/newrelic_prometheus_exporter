package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var detailsLabelName = []string{
	"call_count",
	"total_call_time",
	"total_exclusive_call_time",
	"min_call_time",
	"max_call_time",
	"sum_of_squares",
}
var metricNamespace = "newrelic"

func preConnect(res http.ResponseWriter, req *http.Request) bool {
	log.Info("Handling pre_connect")
	responseData := map[string]string{
		"return_value": "newrelic_prometheus_exporter",
	}

	responseJson, err := json.Marshal(responseData)
	if err != nil {
		panic(err)
	}

	res.Write(responseJson)
	return false
}

func (c *Exporter) handleMetricData(res http.ResponseWriter, req *http.Request) (results bool) {
	results = true
	requestBody, err := parseRequestBodyAsJsonArray(req)

	if err == nil {
		go c.exportNewrelicMetricData(requestBody)
	} else {
		log.Printf("Error reading metric data: %v", err)
	}
	return
}

func (c *Exporter) exportNewrelicMetricData(metricData []JSON) {
	// for more info on how the metric data is structured have a look at this link
	// https://www.evernote.com/shard/s496/sh/82da8165-9fbd-4412-9e79-d70ea70dce5a/081a9b844ccc0cd62139544400b5ccee

	// metrics start from index 3.
	// index 0 is the token sent by agent.
	// index 1 and 2 are start and end time when the data was sent by agent
	metrics := metricData[3].([]interface{})

	// we loop all metrics and add them to our prometheus container
	for i := range metrics {
		// metrics are in list
		// we convert it to a list of interface for easy manipulation
		metric := metrics[i].([]interface{})

		// each metric has two values
		// metric details and metric data
		metricMetaData := metric[0].(map[string]interface{})

		// metric actual data
		metricData := metric[1].([]interface{})

		// get the metric name
		metricName := metricMetaData["name"].(string)
		//
		scope := metricMetaData["scope"].(string)

		if scope == "" {
			scope = "null"
		}

		// labels to use in prometheus
		labels := mergeLabels(map[string]string{
			"scope":       scope,
			"metric_name": metricName,
		}, Config.getGlobalLabel())

		// each metric contains multiple metrics
		// to see all metrics in detailsLabelName
		for i := range metricData {
			// get the metric type based on index
			labels["metric_type"] = detailsLabelName[i]

			// adding metric to be added scrapped
			histogram, err := c.Histograms.Get(metricNamespace, labels, "Newrelic metrics (histogram)")
			if err == nil {
				histogram.Observe(metricData[i].(float64))
			} else {
				// Todo this should be sent to sentry
				fmt.Println("original meta data:", metricData)
				fmt.Println("err:", err)
			}
		}
	}
}

func analyticData(res http.ResponseWriter, req *http.Request) bool {
	return true
}
