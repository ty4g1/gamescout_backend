package models

type User struct {
	ID               string
	SwipeHistory     []int
	PreferenceVector []float64
}
