package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Store the metrics in memory
var (
	metrics = make(map[string]float64)
)

// Example payload:
// {
//     "unit": "1",
//     "metric": "DIO0",
//     "value": 123.45
// }

// Assuming timestamp is Unix time
type ModbusRecord struct {
	Timestamp  int64      `json:"timestamp"`
	Date       string     `json:"date"`
	Bdate      *time.Time `json:"bdate"`
	ServerID   int        `json:"server_id"`
	BServerID  int        `json:"bserver_id"`
	Addr       int        `json:"addr"`
	BAddr      *int       `json:"baddr"`
	FullAddr   string     `json:"full_addr"`
	Size       int        `json:"size"`
	Data       string     `json:"data"`
	RawData    *string    `json:"raw_data"`
	ServerName string     `json:"server_name"`
	IP         string     `json:"ip"`
	Name       string     `json:"name"`
}

type Metric struct {
	Unit   string  `json:"unit"`
	Metric string  `json:"metric"`
	Value  float64 `json:"value"`
}

type ModbusPayload struct {
	Modbus ModbusRecord `json:"Modbus"`
}

func main() {

	http.HandleFunc("/", handleMetricUpdate)
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":3000", nil)

}
