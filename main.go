package main

import (
	"fmt"
	"time"
)

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func drawBox(width int) string {
	var symbol string
	symbol = "□" // Use an empty box character

	box := ""
	for i := 0; i < width; i++ {
		box += symbol + " "
	}
	return box
}

func main() {
	t := time.Now()
	year := t.Year()
	month := t.Month()
	day := t.Day()
	numOfDays := daysInMonth(year, month)

	var boxes []string
	for i := 0; i < numOfDays; i++ {
		boxes = append(boxes, drawBox(1))
	}
	simbol := "■ "
	boxes[day] = simbol

	for i := 0; i < len(boxes); i++ {
		fmt.Printf("%s ", boxes[i])
	}
	fmt.Println()
}
