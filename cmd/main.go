package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"

	"go.fergus.london/tcx/pkg/formats"
)

func main() {
	hike, err := formats.ParseHikeFromTCX("example.tcx")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Start Time: %s\n", hike.StartTime)
	fmt.Printf("Distance: %f\n", hike.Distance)
	fmt.Printf("Duration: %s\n", hike.Duration)
	fmt.Printf("Calories: %d\n", hike.Calories)
	fmt.Printf("MinHR: %d\n", hike.MinHR)
	fmt.Printf("MaxHR: %d\n", hike.MaxHR)

	nTrackpoints := len(hike.Points)
	fmt.Printf("Trackpoints: %d\n", nTrackpoints)

	for _, idx := range randomSample(25, nTrackpoints) {
		fmt.Printf("Point %d: %fs - [%f, %f, %f] - %dbpm\n",
			idx,
			hike.Points[idx].Elapsed.Seconds(),
			hike.Points[idx].Latitude,
			hike.Points[idx].Longitude,
			hike.Points[idx].Altitude,
			hike.Points[idx].HeartRate,
		)
	}

	geojsonFile, err := os.Create("output.geojson.json")
	if err != nil {
		panic(err)
	}
	defer geojsonFile.Close()

	metadataFile, err := os.Create("output.metadata.json")
	if err != nil {
		panic(err)
	}
	defer metadataFile.Close()

	if err := hike.ToGeoJSON(geojsonFile, metadataFile); err != nil {
		panic(err)
	}
}

func randomSample(n int, max int) []int {
	indices := make([]int, n)

	for i := 0; i < n; i++ {
		indices[i] = rand.Intn(max)
	}

	sort.Ints(indices)
	return indices
}
