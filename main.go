package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

func setUptime(status float64, unit string) {
	fmt.Printf("Setting uptime for unit: %s", unit)
	uniqueUpKey := fmt.Sprintf("%s{unit=\"%s\"}", "UP", unit)
	metrics[uniqueUpKey] = status
}

func handleUptime(ticker *time.Ticker, lastUpdate *time.Time) {

	for {
		<-ticker.C
		fmt.Println(lastUpdate)
		// if the last update happened within the last 5 minutes, its still up.
		currentTime := time.Now()

		if currentTime.Add(-1*time.Minute).Unix() > lastUpdate.Unix() {
			fmt.Println("Last update was longer than 5 min ago, assume it is down")
			setUptime(0, "1") // Hardcoded Unit = 1 for now
		}
	}
}

func main() {

	var lastUpdate time.Time

	ticker := time.NewTicker(1 * time.Second)
	app := fiber.New()

	// Handle incoming metrics
	app.Post("/", func(c *fiber.Ctx) error {
		modbusPayload := new(ModbusPayload)
		metric := new(Metric)

		if err := c.BodyParser(modbusPayload); err != nil {
			return err
		}

		fmt.Println(modbusPayload.Modbus)

		metric.Unit = strconv.Itoa(modbusPayload.Modbus.ServerID)
		metric.Metric = modbusPayload.Modbus.Name

		switch {
		case strings.HasPrefix(modbusPayload.Modbus.Name, "DI"):
			// Handle DI type data
			intValue, err := strconv.ParseInt(strings.Trim(modbusPayload.Modbus.Data, "[]"), 10, 64)
			if err != nil || (intValue != 0 && intValue != 1) {
				log.Printf("Error or invalid value for DI data type: %v", err)
				intValue = 0
			}
			metric.Value = float64(intValue)

		case strings.HasPrefix(modbusPayload.Modbus.Name, "ANA"):
			// Handle ANA type data
			floatValue, err := strconv.ParseFloat(strings.Trim(modbusPayload.Modbus.Data, "[]"), 64)
			if err != nil {
				log.Printf("Error parsing ANA data to float: %v", err)
				floatValue = 0.0
			}
			metric.Value = floatValue

		case strings.HasPrefix(modbusPayload.Modbus.Name, "DO"):
			// Handle D0 type data as string
			metric.Value = float64(len(modbusPayload.Modbus.Data))

		default:
			// Default case
			metric.Value = 0
		}

		// // Construct a unique key with unit and metric
		fmt.Printf("Unit: %s", metric.Unit)
		uniqueKey := fmt.Sprintf("%s{unit=\"%s\"}", metric.Metric, metric.Unit)
		if metric.Unit != "0" {
			metrics[uniqueKey] = metric.Value
			fmt.Printf("Metric updated: %s = %f\n", uniqueKey, metric.Value)
		}

		setUptime(1, metric.Unit)

		lastUpdate = time.Now()

		return nil
	})

	// Expose the metrics for Prometheus
	app.Get("/metrics", func(c *fiber.Ctx) error {
		response := ""
		for key, value := range metrics {
			response += fmt.Sprintf("%s %f\n", key, value)
		}
		return c.SendString(response)
	})

	go handleUptime(ticker, &lastUpdate)

	log.Fatal(app.Listen(":3000"))
}
