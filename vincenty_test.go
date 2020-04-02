package vincenty

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"testing"
	"os"
	"io"
	"fmt"
	"strconv"
	"time"
	"encoding/csv"
)

func verifyInverse(p1, p2 LatLng, expectedResult float64, t *testing.T) {
	dist := Inverse(p1, p2)
	if dist.Metres() != expectedResult {
		t.Errorf("Inverse() returned %v -- expected distance: %v", dist.Metres(), expectedResult)
	}
}

func TestCoincident(t *testing.T) {
	point1 := LatLng{Latitude: 0.0, Longitude: 0.0}
	point2 := LatLng{Latitude: 0.0, Longitude: 0.0}
	expectedDist := 0.0

	verifyInverse(point1, point2, expectedDist, t)
}

func Test0001(t *testing.T) {
	point1 := LatLng{Latitude: 0.0, Longitude: 0.0}
	point2 := LatLng{Latitude: 0.0, Longitude: 1.0}
	expectedDist := 111319.491

	verifyInverse(point1, point2, expectedDist, t)
}

func Test0010(t *testing.T) {
	point1 := LatLng{Latitude: 0.0, Longitude: 0.0}
	point2 := LatLng{Latitude: 1.0, Longitude: 0.0}
	expectedDist := 110574.389

	verifyInverse(point1, point2, expectedDist, t)
}

func TestSlowConvergence(t *testing.T) {
	point1 := LatLng{Latitude: 0.0, Longitude: 0.0}
	point2 := LatLng{Latitude: 0.5, Longitude: 179.5}
	expectedDist := 19936288.579

	verifyInverse(point1, point2, expectedDist, t)
}

func TestFailureToConverge(t *testing.T) {
	point1 := LatLng{Latitude: 0.0, Longitude: 0.0}
	point2 := LatLng{Latitude: 0.5, Longitude: 179.7}
	expectedDist := -1.0

	verifyInverse(point1, point2, expectedDist, t)
}

func TestBostonNewYork(t *testing.T) {
	Boston := LatLng{Latitude: 42.3541165, Longitude: -71.0693514}
        NewYork := LatLng{Latitude: 40.7791472, Longitude: -73.9680804}
	expectedDist := 298396.057

	verifyInverse(Boston, NewYork, expectedDist, t)
}

func TestSpeed(t *testing.T) {
	var testData [][]float64
	testDataFile, err := os.Open("testdata.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(testDataFile)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fRow := []float64{}
		for _, s := range(row) {
			if f, err := strconv.ParseFloat(s, 64); err != nil {
				panic(err)
			} else {
				fRow = append(fRow, f)
			}
		}
		testData = append(testData, fRow)
	}
	testDataFile.Close()
	before := time.Now()
	for _, data := range(testData) {
		dist := Inverse(LatLng{Latitude: data[0], Longitude: data[1]}, LatLng{Latitude: data[2], Longitude: data[3]})
		if dist.Metres() - data[4] > 0.0011 {
			t.Errorf("Inverse() returned %v -- expected distance: %v", dist.Metres(), data[4])
		}
	}
	fmt.Printf("Timing (%v tests): %vs\n", len(testData), time.Since(before).Seconds())
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
