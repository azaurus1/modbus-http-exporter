package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handleMetricUpdate(w http.ResponseWriter, req *http.Request) {
	metric := new(Metric)

	// unmarshall request to modbus payload

	var modbusPayload ModbusPayload

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	json.Unmarshal(reqBody, &modbusPayload)

	fmt.Println(modbusPayload.Modbus)

	metric.Unit = strconv.Itoa(modbusPayload.Modbus.ServerID)
	metric.Metric = modbusPayload.Modbus.Name

	switch {
	case strings.HasPrefix(modbusPayload.Modbus.Name, "DI"):
		log.Println("Digital input update")

		// Handle DI type data
		intValue, err := strconv.ParseInt(strings.Trim(modbusPayload.Modbus.Data, "[]"), 10, 64)
		if err != nil || (intValue != 0 && intValue != 1) {
			log.Printf("Error or invalid value for DI data type: %v", err)
			intValue = 0
		}

		di1.Set(float64(intValue))
		log.Printf("DI1 set to: %v", float64(intValue))

	case strings.HasPrefix(modbusPayload.Modbus.Name, "ANA"):

		log.Println("Analog input update")

		// Handle ANA type data
		floatValue, err := strconv.ParseFloat(strings.Trim(modbusPayload.Modbus.Data, "[]"), 64)
		if err != nil {
			log.Printf("Error parsing ANA data to float: %v", err)
			floatValue = 0.0
		}

		ana1.Set(floatValue)
		log.Printf("ANA1 set to: %v", floatValue)

	case strings.HasPrefix(modbusPayload.Modbus.Name, "DO"):
		// Handle D0 type data as string
		metric.Value = float64(len(modbusPayload.Modbus.Data))

	default:
		// Default case
		metric.Value = 0
	}
}
