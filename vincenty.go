package vincenty

// Go re-implementation of https://github.com/maurycyp/vincenty

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"math"
)

// Distance is the return type of vincenty.Inverse()
// Use Metres() or Kilometres() to get the distance in the unit of your choice.
// If you prefer imperial units, use NauticalMiles() Miles() or Feet().
// Or if you're in the US but like SI standards, you may want to use Meters() or Kilometers()  :)
type Distance float64

// Metres returns the Distance d in metres
func (d Distance)Metres() float64 {
	return float64(d)
}

// Meters also returns the Distance d in metres, but in US English
func (d Distance)Meters() float64 {
	return float64(d)
}

// Kilometres returns the Distance d in kilometres
func (d Distance)Kilometres() float64 {
	return float64(d) / 1000.0
}

// Kilometers also returns the Distance d in kilometres, but in US English
func (d Distance)Kilometers() float64 {
	return float64(d) / 1000.0
}

// NauticalMiles returns the Distance d in, you guessed it, nautical miles
func (d Distance)NauticalMiles() float64 {
	return float64(d) / 1852.0
}

// Miles returns the Distance d in miles
func (d Distance)Miles() float64 {
	return float64(d) / 1609.344
}

// Feet returns the Distance d in feet
func (d Distance)Feet() float64 {
	return float64(d) / 0.3048
}

// LatLng represents a point on Earth defined by its Latitude and Longitude
type LatLng struct {
	Latitude float64
	Longitude float64
}

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func square(f float64) float64 {
	return f * f
}

// Inverse calculates the distance between two points on the surface of a spheroid
// using Vincenty's formula (inverse method)
func Inverse(point1, point2 LatLng) Distance {
	// WGS 84
	a := 6378137.0  // meters
	f := 1.0 / 298.257223563
	b := 6356752.314245  // meters; b = (1 - f)a

	MaxIterations := 200
	ConvergenceThreshold := 1e-12  // .000,000,000,001

	// short-circuit coincident points
	if point1.Latitude == point2.Latitude && point1.Longitude == point2.Longitude {
		return Distance(0.0)
	}

	U1 := math.Atan((1.0 - f) * math.Tan(radians(point1.Latitude)))
	U2 := math.Atan((1.0 - f) * math.Tan(radians(point2.Latitude)))
	L := radians(point2.Longitude - point1.Longitude)
	Lambda := L

	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)

	for i := 0; i < MaxIterations; i++ {
		sinLambda := math.Sin(Lambda)
		cosLambda := math.Cos(Lambda)
		sinSigma := math.Sqrt(square(cosU2 * sinLambda) + square(cosU1 * sinU2 - sinU1 * cosU2 * cosLambda))
		if sinSigma == 0.0 {
			return Distance(0.0)  // coincident points
		}
		cosSigma := sinU1 * sinU2 + cosU1 * cosU2 * cosLambda
		sigma := math.Atan2(sinSigma, cosSigma)
		sinAlpha := cosU1 * cosU2 * sinLambda / sinSigma
		cosSqAlpha := 1.0 - square(sinAlpha)
		cos2SigmaM := 0.0
		if cosSqAlpha != 0 {
			cos2SigmaM = cosSigma - 2.0 * sinU1 * sinU2 / cosSqAlpha
		}
		C := f / 16.0 * cosSqAlpha * (4.0 + f * (4.0 - 3.0 * cosSqAlpha))
		LambdaPrev := Lambda
		Lambda = L + (1.0 - C) * f * sinAlpha * (sigma + C * sinSigma * (cos2SigmaM + C * cosSigma * (-1.0 + 2.0 * square(cos2SigmaM))))
		if math.Abs(Lambda - LambdaPrev) < ConvergenceThreshold {
			// successful convergence
			uSq := cosSqAlpha * (square(a) - square(b)) / square(b)
			A := 1.0 + uSq / 16384.0 * (4096.0 + uSq * (-768.0 + uSq * (320.0 - 175.0 * uSq)))
			B := uSq / 1024.0 * (256.0 + uSq * (-128.0 + uSq * (74.0 - 47.0 * uSq)))
			deltaSigma := B * sinSigma * (cos2SigmaM + B / 4.0 * (cosSigma * (-1.0 + 2.0 * square(cos2SigmaM)) - B / 6.0 * cos2SigmaM * (-3.0 + 4.0 * square(sinSigma)) * (-3.0 + 4.0 * square(cos2SigmaM))))
			s := b * A * (sigma - deltaSigma)
			s = math.Round(s * 1000)/1000
			return Distance(s)
		}
	}
	return Distance(-1.0)
}
