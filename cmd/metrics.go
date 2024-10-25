package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	di1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "digital_input_1",
		Help: "The value of digital input 1",
	})

	ana1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "analog_input_1",
		Help: "The value of analog input 1",
	})

	do1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "digital_output_1",
		Help: "The value of digital output 1",
	})
)
