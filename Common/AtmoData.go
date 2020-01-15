package common

import "time"

type SensorData struct {
  Timestamp   time.Time
  Temperature float64
  Humidity    float64
  Pressure    float64
}
