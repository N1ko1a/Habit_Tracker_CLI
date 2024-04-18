package main

import (
	"context"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

// Broj ispunjenih navika u danu
func countCompletedInDay(year int, month time.Month, day int) (int, error) {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	//Nalazimo sve navike sa ovom godinom,mesecom i danom koji smo naveli
	filter := bson.D{{Key: "year", Value: year}, {Key: "month", Value: month}, {Key: fmt.Sprintf("boxes.%d", day), Value: "■ "}}

	count, err := Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("Error finding habits completed on day %d: %v", day, err)
	}

	return int(count), nil
}

// Crtanje grafika
func plotingGraph(completedData []int, year int, month time.Month) {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	//Postavljamo labels
	data := completedData
	labels := make([]string, len(data))
	for i := range labels {
		labels[i] = fmt.Sprintf("%d", i+1)
	}

	maxY := 6
	// Create the plot widget
	plot := widgets.NewPlot()
	plot.Title = "Habits Completed"
	plot.SetRect(2, 4, 70, 30)
	plot.Data = make([][]float64, 1)
	plot.Data[0] = make([]float64, len(data))
	for i, val := range data {
		// Normalize the data to fit the range of the y-axis
		if val > maxY {
			plot.Data[0][i] = float64(maxY)
		} else {
			plot.Data[0][i] = float64(val)
		}
	}
	plot.DataLabels = labels

	// Change plot line color to white
	plot.LineColors[0] = ui.ColorWhite

	// Adjust horizontal scale to stretch out lines
	plot.HorizontalScale = 2.0

	fmt.Printf("\n\n  %d %s\n", year, month)
	// Render the plot
	ui.Render(plot)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
