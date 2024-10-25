package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// sends post request with the payload

type ModbusPayload struct {
	Modbus ModbusRecord `json:"Modbus"`
}

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

// strings.HasPrefix(modbusPayload.Modbus.Name, "DI") - is where the metric is determined
// DI, ANA, DO
// intValue, err := strconv.ParseInt(strings.Trim(modbusPayload.Modbus.Data, "[]"), 10, 64) => metric.Value = float64(intValue)

func main() {

	// build payloads
	// send them to server
	//

	// Example date and optional date (bdate)
	date := "2024-06-04"
	bdate := time.Now()

	// Example optional fields
	var baddr int = 2
	var rawData string = "rawDataExample"

	record := ModbusRecord{
		Timestamp:  1717454045,
		Date:       date,
		Bdate:      &bdate,
		ServerID:   1,
		BServerID:  1,
		Addr:       1,
		BAddr:      &baddr,
		FullAddr:   "",
		Size:       16,
		Data:       "[15]",
		RawData:    &rawData,
		ServerName: "mock",
		IP:         "127.0.0.1",
		Name:       "mock",
	}

	payload := map[string]ModbusRecord{"Modbus": record}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("error marshalling payload", err)
	}

	log.Println(string(jsonData))

	// send to server
	r, err := http.NewRequest("POST", "http://localhost:3000/", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("can not send post request to exporter", err)
	}

	log.Println(r)

	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		log.Println("can not send post request to exporter", err)
	}

	log.Println(res)

}
