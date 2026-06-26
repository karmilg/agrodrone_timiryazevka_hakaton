package models

import "time"

type Cell struct {
	X           int
	Y           int
	Soil        float64
	NDVI        float64
	Temp        float64
	LastScanned time.Time
	Scanned     bool
}

func (c *Cell) IsProblem() bool {
	return c.NDVI < 0.3 || c.Soil < 25
}