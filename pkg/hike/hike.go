package hike

import (
	"encoding/json"
	"time"
)

// Hike represents a hike, including the start time, distance, duration,
// calories, minimum and maximum heart rate, and the points (latitude,
// longitude, altitude, elapsed time, and heart rate) which make up the
// route.
type Hike struct {
	StartTime    time.Time
	Distance     float64
	Duration     time.Duration
	Calories     int
	MinHR, MaxHR int
	Points       []Point
}

// Point represents a point on the route, including the latitude, longitude,
// and altitude, as well as the elapsed time and heart rate associated with
// the hiker.
type Point struct {
	Latitude, Longitude, Altitude float64
	Elapsed                       time.Duration `json:"secs"`
	HeartRate                     int           `json:"hr"`
}

// MarshalJSON implements the json.Marshaler interface for the Point type,
// due to how we want to output the "metadata" associated with a Point (i.e.
// non-geographic data, such as heart rate and elapsed time), we simply
// marshal those values ONLY.
func (p *Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"hr":   p.HeartRate,
		"secs": p.Elapsed.Seconds(),
	})
}
