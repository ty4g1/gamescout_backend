package utils

import (
	"errors"
	"math"
)

func AddVectors(v1 []float64, v2 []float64) ([]float64, error) {
	if len(v1) == 0 {
		return v2, nil
	}
	if len(v2) == 0 {
		return v1, nil
	}

	if len(v1) != len(v2) {
		return nil, errors.New("can only sum same size vectors")
	}
	n := len(v1)
	res := make([]float64, n)
	for i := range n {
		res[i] = v1[i] + v2[i]
	}
	return res, nil
}

func SubtractVectors(v1 []float64, v2 []float64) ([]float64, error) {
	if len(v1) == 0 {
		return v2, nil
	}
	if len(v2) == 0 {
		return v1, nil
	}

	if len(v1) != len(v2) {
		return nil, errors.New("can only subtract same size vectors")
	}
	n := len(v1)
	res := make([]float64, n)
	for i := range n {
		res[i] = v1[i] - v2[i]
	}
	return res, nil
}

func SumRows(v [][]float64) ([]float64, error) {
	if len(v) == 0 {
		return []float64{}, nil
	}
	n := len(v[0])
	res := make([]float64, n)
	for _, vector := range v {
		sum, err := AddVectors(res, vector)
		if err != nil {
			return nil, err
		}
		res = sum
	}
	return res, nil
}

// In your utils package
func NormalizeVector(v []float64) []float64 {
	var magnitude float64
	for _, val := range v {
		magnitude += val * val
	}
	magnitude = math.Sqrt(magnitude)

	if magnitude == 0 {
		return v
	}

	normalized := make([]float64, len(v))
	for i, val := range v {
		normalized[i] = val / magnitude
	}

	return normalized
}

func ComputeSimilarity(v1 []float64, v2 []float64) (float64, error) {
	if len(v1) == 0 || len(v2) == 0 {
		return 0, nil
	}

	if len(v1) != len(v2) {
		return 0, errors.New("can only compute similarity for same size vectors")
	}
	n := len(v1)
	var res float64 = 0
	for i := range n {
		res += v1[i] * v2[i]
	}
	return res, nil
}
