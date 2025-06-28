package main

import (
	"os"
)

const (
	kafkaBroker = "kafka:9092"
	topic       = "sensorData"
)

func main() {

	return
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
