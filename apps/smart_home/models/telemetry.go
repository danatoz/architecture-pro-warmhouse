package models

import "time"

type SensorData struct {
	SensorID  int       `json:"sensor_id"`
	Type      string    `json:"type"`
	Unit      string    `json:"unit"`
	Value     *float64  `json:"value"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
