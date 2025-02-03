package formats

import (
	"encoding/xml"
	"os"
	"time"

	"go.fergus.london/tcx/pkg/hike"
)

// TrainingCenterDatabase represents the root XML element for a FitBit-generated
// TCX file. It contains the mapping for the key XML values that we wish to parse
// in to a Hike struct.
type TrainingCenterDatabase struct {
	Activities struct {
		Activity []struct {
			Time string `xml:"Id"` // Gotcha: The `Id` element is actually the start time of the activity
			Lap  []struct {
				StartTime        string  `xml:"StartTime,attr"`
				TotalTimeSeconds float64 `xml:"TotalTimeSeconds"`
				DistanceMeters   float64 `xml:"DistanceMeters"`
				Calories         int     `xml:"Calories"`
				Track            []struct {
					Time      string  `xml:"Time"`
					Latitude  float64 `xml:"Position>LatitudeDegrees"`
					Longitude float64 `xml:"Position>LongitudeDegrees"`
					Altitude  float64 `xml:"AltitudeMeters"`
					HeartRate int     `xml:"HeartRateBpm>Value"`
				} `xml:"Track>Trackpoint"`
			} `xml:"Lap"`
		} `xml:"Activity"`
	} `xml:"Activities"`
}

// ParseHikeFromTCX reads a TCX file and returns a fully populated Hike struct,
// complete with all points and stats - _it does not perform any sampling and
// will likely return a high number of points_!
func ParseHikeFromTCX(path string) (*hike.Hike, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var db TrainingCenterDatabase
	if err := xml.Unmarshal(data, &db); err != nil {
		return nil, err
	}

	h := hike.Hike{}

	for _, activity := range db.Activities.Activity {
		startTime, err := time.Parse(time.RFC3339, activity.Time)
		if err != nil {
			return nil, err
		}

		totalPoints := 0
		for _, lap := range activity.Lap {
			totalPoints += len(lap.Track)
		}

		h.Points = make([]hike.Point, totalPoints)
		pointIdx := 0

		for _, lap := range activity.Lap {
			h.Distance += lap.DistanceMeters
			h.Duration += time.Duration(lap.TotalTimeSeconds * float64(time.Second))
			h.Calories += lap.Calories

			for _, trackpoint := range lap.Track {
				if minHr := h.MinHR; minHr == 0 || minHr > trackpoint.HeartRate {
					h.MinHR = trackpoint.HeartRate
				}
				if trackpoint.HeartRate > h.MaxHR {
					h.MaxHR = trackpoint.HeartRate
				}

				elapsed, err := time.Parse(time.RFC3339, trackpoint.Time)
				if err != nil {
					return nil, err
				}

				h.Points[pointIdx] = hike.Point{
					Elapsed:   elapsed.Sub(startTime),
					Latitude:  trackpoint.Latitude,
					Longitude: trackpoint.Longitude,
					Altitude:  trackpoint.Altitude,
					HeartRate: trackpoint.HeartRate,
				}

				pointIdx++
			}
		}
	}

	return &h, nil
}
