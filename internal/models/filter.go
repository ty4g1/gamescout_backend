package models

import "time"

type PriceRange struct {
	Min int
	Max int
}

type ReleaseDate struct {
	IsBefore bool
	Date     time.Time
}
