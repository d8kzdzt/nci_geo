package nci_geo_lib

import (
	"fmt"
	"math"
)

func checkNotWGS(source, dest *SR) bool {
	return ((source.datum.datum_type == pjd3Param || source.datum.datum_type == pjd7Param) && dest.DatumCode != "WGS84")
}

const enu = "enu"
const longlat = "longlat"

// CoordinateTransform creates a function that transforms a point from sr
// to the destination spatial reference.
func (source *SR) CoordinateTransform(dest *SR) (Transformer, error) {
	if dest == nil {
		return nil, fmt.Errorf("proj: destination is nil")
	}

	// If source and dest are the same, we don't need to do any transforming
	const ulpTolerance = 3 // Our tolerance is 3 units in the last place
	if source.Equal(dest, 3) {
		return func(x, y float64) (float64, float64, error) {
			return x, y, nil
		}, nil
	}

	return func(x, y float64) (float64, float64, error) {
		point := []float64{x, y}
		_, sourceInverse, err := source.Transformers()
		if err != nil {
			return math.NaN(), math.NaN(), err
		}
		destForward, _, err := dest.Transformers()
		if err != nil {
			return math.NaN(), math.NaN(), err
		}

		// DGR, 2010/11/12
		if source.Axis != enu {
			point, err = adjust_axis(source, false, point)
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
		}
		// Transform source points to long/lat, if they aren't already.
		if source.Name == longlat {
			point[0] *= deg2rad // convert degrees to radians
			point[1] *= deg2rad
		} else {
			point[0] *= source.ToMeter
			point[1] *= source.ToMeter
			point[0], point[1], err = sourceInverse(point[0], point[1]) // Convert Cartesian to longlat
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
		}
		// Adjust for the prime meridian if necessary
		if !math.IsNaN(source.FromGreenwich) {
			point[0] += source.FromGreenwich
		}

		// Convert datums if needed, and if possible.
		z := 0.
		if len(point) == 3 {
			z = point[2]
		}
		point[0], point[1], z, err = datumTransform(source.datum, dest.datum,
			point[0], point[1], z)
		if err != nil {
			return math.NaN(), math.NaN(), err
		}
		if len(point) == 3 {
			point[2] = z
		}

		// Adjust for the prime meridian if necessary
		if !math.IsNaN(dest.FromGreenwich) {
			point[0] -= dest.FromGreenwich
		}

		if dest.Name == longlat {
			// convert radians to decimal degrees
			point[0] *= r2d
			point[1] *= r2d
		} else { // else project
			point[0], point[1], err = destForward(point[0], point[1])
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
			point[0] /= dest.ToMeter
			point[1] /= dest.ToMeter
		}

		// DGR, 2010/11/12
		if dest.Axis != enu {
			point, err = adjust_axis(dest, true, point)
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
		}

		return point[0], point[1], nil
	}, nil
}

// CoordinateTransform creates a function that transforms a point from sr
// to the destination spatial reference.
func (source *SR) CoordinateTransformNew(dest *SR) (Transformer, error) {
	if dest == nil {
		return nil, fmt.Errorf("proj: destination is nil")
	}

	// If source and dest are the same, we don't need to do any transforming
	const ulpTolerance = 3 // Our tolerance is 3 units in the last place
	if source.Equal(dest, 3) {
		return func(x, y float64) (float64, float64, error) {
			return x, y, nil
		}, nil
	}
	//
	_, sourceInverse, err := source.Transformers()
	if err != nil {
		return func(x, y float64) (float64, float64, error) {
			return math.NaN(), math.NaN(), err
		}, nil
	}
	destForward, _, err := dest.Transformers()
	if err != nil {
		return func(x, y float64) (float64, float64, error) {
			return math.NaN(), math.NaN(), err
		}, nil
	}
	return func(x, y float64) (float64, float64, error) {
		point := []float64{x, y}
		// DGR, 2010/11/12
		if source.Axis != enu {
			point, err = adjust_axis(source, false, point)
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
		}
		// Transform source points to long/lat, if they aren't already.
		if source.Name == longlat {
			point[0] *= deg2rad // convert degrees to radians
			point[1] *= deg2rad
		} else {
			point[0] *= source.ToMeter
			point[1] *= source.ToMeter
			point[0], point[1], err = sourceInverse(point[0], point[1]) // Convert Cartesian to longlat
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
		}
		// Adjust for the prime meridian if necessary
		if !math.IsNaN(source.FromGreenwich) {
			point[0] += source.FromGreenwich
		}

		// Convert datums if needed, and if possible.
		z := 0.
		if len(point) == 3 {
			z = point[2]
		}
		point[0], point[1], z, err = datumTransform(source.datum, dest.datum,
			point[0], point[1], z)
		if err != nil {
			return math.NaN(), math.NaN(), err
		}
		if len(point) == 3 {
			point[2] = z
		}

		// Adjust for the prime meridian if necessary
		if !math.IsNaN(dest.FromGreenwich) {
			point[0] -= dest.FromGreenwich
		}

		if dest.Name == longlat {
			// convert radians to decimal degrees
			point[0] *= r2d
			point[1] *= r2d
		} else { // else project
			point[0], point[1], err = destForward(point[0], point[1])
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
			point[0] /= dest.ToMeter
			point[1] /= dest.ToMeter
		}

		// DGR, 2010/11/12
		if dest.Axis != enu {
			point, err = adjust_axis(dest, true, point)
			if err != nil {
				return math.NaN(), math.NaN(), err
			}
		}

		return point[0], point[1], nil
	}, nil
}
