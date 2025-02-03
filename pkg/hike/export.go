package hike

import (
	"encoding/json"
	"io"

	"github.com/everystreet/go-geojson/v3"
)

// Exports the hike to a file containing the route (incl. elevation) expressed
// as a GeoJSON LineString, and - optionally - a file containing the metadata
// (heart rate and elapsed time) expressed as a JSON array.
func (h *Hike) ToGeoJSON(w, metadata io.Writer) error {
	coords := make([]geojson.Position, len(h.Points))
	for i, p := range h.Points {
		coords[i] = geojson.MakePositionWithElevation(p.Latitude, p.Longitude, p.Altitude)
	}

	feature := geojson.NewFeature(
		geojson.NewLineString(coords[0], coords[1], coords[2:]...),
		geojson.Property{Name: "distance", Value: h.Distance},
		geojson.Property{Name: "duration", Value: h.Duration.Seconds()},
		geojson.Property{Name: "calories", Value: h.Calories},
		geojson.Property{Name: "minHR", Value: h.MinHR},
		geojson.Property{Name: "maxHR", Value: h.MaxHR},
	)

	// Encode and write the GeoJSON
	var err error
	if err = json.NewEncoder(w).Encode(feature); err == nil {
		if metadata != nil {
			err = json.NewEncoder(metadata).Encode(h.Points)
		}
	}

	return err
}
