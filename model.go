package main

import (
	"time"
)

type Habit struct {
	Name  string     `json:"name"`
	Boxes []string   `json:"boxes"`
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
}
