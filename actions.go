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
	if err != nil{
		panic(err)
	}

	res.Write(responseJson)
	return false
}



func (c * Exporter) handleMetricData(res http.ResponseWriter, req *http.Request) (results bool){
	results = true
	log.Info("Handling metric_data")
	requestBody, err := parseRequestBodyAsJsonArray(req)

	if err == nil{
		go c.exportNewrelicMetricData(requestBody)
	} else {
		log.Printf("Error reading metric data: %v", err)
	}
	return
}

func (c *Exporter) exportNewrelicMetricData(metricData []JSON){
	metrics := metricData[3].([]interface{})
	for i := range metrics{
		metric := metrics[i].([]interface{})
		metaData := metric[0].(map[string]interface{})
		details := metric[1].([]interface{})
		rawName := metaData["name"].(string)
		metricName := metricNamespace
		scope := metaData["scope"].(string)
		for i := range details{
			labels := map[string]string{
				"scope": scope,
				"metric_name": rawName,
			}
			labels["metric_type"] = detailsLabelName[i]
			gauge, err := c.Gauges.Get(metricName, labels, "Newrelic metrics")
			if err==nil{
				gauge.Set(details[i].(float64))
			}else{
				fmt.Println("original meta data:", metaData)
				fmt.Println("err:", err)
			}
		}
	}
}

func analyticData(res http.ResponseWriter, req *http.Request) bool{
	log.Info("Handling analytic_event_data")
	requestBody, err := parseRequestBody(req)
	if err != nil{
		log.Printf("Error reading analytic_event_data: %v", err)
	}else{
		log.Info("analytic_event_data:", requestBody)
	}
	return true
}