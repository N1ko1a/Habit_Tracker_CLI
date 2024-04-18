package main

import (
	"fmt"
	"time"
)

// Uzimanje broj dana u mesecu
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day() //Prosledjujemo index unutar niza boxes pa posto on krece od 0 moramo da dodamo +1
}

// Crtanje kockica posot je kocka strign vracamo string
func drawBox(width int) string {
	var symbol string
	symbol = "□"

	box := ""
	for i := 0; i < width; i++ {
		box += symbol + " "
	}
	return box
}

// Proveravamo da li je index koji stavljamo za dan validan
func isValidDay(year int, month time.Month, day int) bool {
	daysInMonth := daysInMonth(year, month)
	return day >= 1 && day <= daysInMonth // proveri da li se dan nalazi izmedju 1 i daysInMonth (31)
}

// Parsamo string u time.Month
func parseMonth(monthString string) (time.Month, error) {

	// List of months
	months := map[string]time.Month{
		"January":   time.January,
		"February":  time.February,
		"March":     time.March,
		"April":     time.April,
		"May":       time.May,
		"June":      time.June,
		"July":      time.July,
		"August":    time.August,
		"September": time.September,
		"October":   time.October,
		"November":  time.November,
		"December":  time.December,
	}
	//Nadji mesec u mapi
	month, ok := months[monthString]
	if !ok {
		return 0, fmt.Errorf("Invalid month: %s", monthString)
	}
	return month, nil
}
