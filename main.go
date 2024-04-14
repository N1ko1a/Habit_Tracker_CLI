package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Habit struct {
	Name  string     `json:"name"`
	Boxes []string   `json:"boxes"`
	Year  int        `json:"year"`
	Month time.Month `json:"month"`
}

// kreiramo instancu da bi se povezali sa bazom i imali interakcije
var mongoClient *mongo.Client

const uri = "mongodb://localhost:27017"

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

func addHabit(habit Habit) error {
	ctx := context.TODO()                            //vraca prazan context koji sluzi kao placeholder
	Database := mongoClient.Database("HabitTracker") //Ovo je baza, koristimo instancu inzad koju smo kreirali
	Collection := Database.Collection("Habits")      //Ovo je colekcija, koristimo umesto instance ime baze

	//Proverava da li postoji vec user sa istim email-om
	//Kreira se filter koji specifira da trazimo email vrednosti user.Email
	filter := bson.D{{Key: "name", Value: habit.Name}, {Key: "year", Value: habit.Year}, {Key: "month", Value: habit.Month}}
	existingDoc := Collection.FindOne(ctx, filter) //koristimo findOne i kao argument dodajemo filter iznad sto smo kreirali
	// U slucaju da ne postoji email vraca nil a ako postoji vraca vrednost tog email

	//mogli smo ovo Err da izbacimo
	if existingDoc.Err() == nil {
		// Document with the same email already exists, skip insertion
		fmt.Printf("Skipping insertion for habit - Name: %s  as it already exists in the %d %d.\n", habit.Name, habit.Year, habit.Month)
		return nil
		// Ako izbacimo Err bilo bi else{}
	} else if existingDoc.Err() != mongo.ErrNoDocuments {
		// An error occurred while checking for existing documents
		return existingDoc.Err()
	}
	_, err := Collection.InsertOne(ctx, bson.D{
		{Key: "name", Value: habit.Name},
		{Key: "boxes", Value: habit.Boxes},
		{Key: "year", Value: habit.Year},
		{Key: "month", Value: habit.Month},
	})

	return err

}

func getAllHabits() {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")
	sortOptions := bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}

	cursor, err := Collection.Find(ctx, bson.D{}, options.Find().SetSort(sortOptions))
	if err != nil {
		fmt.Println("Error querying habits from the database:", err)
		return
	}
	defer cursor.Close(ctx)

	// Process retrieved data from the cursor as needed.
	var habits []Habit
	if err := cursor.All(ctx, &habits); err != nil {
		fmt.Println("Error decoding habits from the cursor:", err)
		return
	}

	// Print the retrieved habits
	var prevMonth time.Month
	var prevYear int

	for _, habit := range habits {
		if habit.Month != prevMonth || habit.Year != prevYear {
			fmt.Println()
			fmt.Println()
			fmt.Printf("%s %d\n", habit.Month, habit.Year) // Add newline if month differs
			fmt.Println()
		}
		// Define the width of the name field
		nameWidth := 15 // Change this to the desired width

		// Print the name field with defined width
		fmt.Printf("%-*s", nameWidth, habit.Name)

		for i := 0; i < len(habit.Boxes); i++ {
			fmt.Printf("%s ", habit.Boxes[i])
		}
		fmt.Println()

		prevMonth = habit.Month // Update prevMonth for the next iteration
		prevYear = habit.Year   // Update prevMonth for the next iteration
		// Depending on your data structure, you might need to adjust the print format
	}
	fmt.Println()
	fmt.Println()
}

func getCurrentHabits(year int, month time.Month) {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")
	sortOptions := bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}

	filter := bson.M{}
	filter["year"] = year
	filter["month"] = month
	cursor, err := Collection.Find(ctx, filter, options.Find().SetSort(sortOptions))
	if err != nil {
		fmt.Println("Error querying habits from the database:", err)
		return
	}
	defer cursor.Close(ctx)

	// Process retrieved data from the cursor as needed.
	var habits []Habit
	if err := cursor.All(ctx, &habits); err != nil {
		fmt.Println("Error decoding habits from the cursor:", err)
		return
	}

	// Print the retrieved habits
	var prevMonth time.Month
	var prevYear int

	for _, habit := range habits {
		if habit.Month != prevMonth || habit.Year != prevYear {
			fmt.Println()
			fmt.Println()
			fmt.Printf("%s %d\n", habit.Month, habit.Year) // Add newline if month differs
			fmt.Println()
		}
		// Define the width of the name field
		nameWidth := 15 // Change this to the desired width

		// Print the name field with defined width
		fmt.Printf("%-*s", nameWidth, habit.Name)

		for i := 0; i < len(habit.Boxes); i++ {
			fmt.Printf("%s ", habit.Boxes[i])
		}
		fmt.Println()

		prevMonth = habit.Month // Update prevMonth for the next iteration
		prevYear = habit.Year   // Update prevMonth for the next iteration
		// Depending on your data structure, you might need to adjust the print format
	}
	fmt.Println()
	fmt.Println()
}

func editHabitName(name string, year int, month time.Month, newName string) error {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Define the filter to find the habit to be edited
	filter := bson.D{{Key: "name", Value: name}, {Key: "year", Value: year}, {Key: "month", Value: month}}

	// Define the update to set the new name
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: newName}}}}

	// Perform the update operation
	_, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating habit name: %v", err)
	}

	return nil
}

func deleteHabit(name string, year int, month time.Month) error {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Define the filter to find the habit to be deleted
	filter := bson.D{{Key: "name", Value: name}, {Key: "year", Value: year}, {Key: "month", Value: month}}

	// Perform the deletion operation
	_, err := Collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting habit: %v", err)
	}

	return nil
}

// Function to parse month string into time.Month
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

	// Lookup the month in the map
	month, ok := months[monthString]
	if !ok {
		return 0, fmt.Errorf("invalid month: %s", monthString)
	}

	return month, nil
}

func compleatedHabit(name string, year int, month time.Month, day int) error {

	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Define the filter to find the habit to be edited
	filter := bson.D{{Key: "name", Value: name}, {Key: "year", Value: year}, {Key: "month", Value: month}}

	simbol := "■ "
	// Define the update to set the new name
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("boxes.%d", day), Value: simbol}}}}

	// Perform the update operation
	_, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating habit name: %v", err)
	}

	return nil
}

func timesCompletedInMonth(name string, year int, month time.Month) (int, error) {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Define the filter to find the habit
	filter := bson.M{
		"name":  name,
		"year":  year,
		"month": month,
		"boxes": "■ ", // Filter for documents where the "boxes" array contains the symbol "■"
	}

	// Perform the aggregation to count the occurrences of "■" in the "boxes" array
	cursor, err := Collection.Aggregate(ctx, mongo.Pipeline{
		{{"$match", filter}},
		{{"$project", bson.M{"_id": 0, "boxes": 1}}},
		{{"$unwind", "$boxes"}},
		{{"$match", bson.M{"boxes": "■ "}}},
		{{"$group", bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}},
	})
	if err != nil {
		return 0, fmt.Errorf("error counting completed habits: %v", err)
	}
	defer cursor.Close(ctx)

	// Extract the count from the result
	var result struct {
		Count int `bson:"count"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("error decoding result: %v", err)
		}
	}

	fmt.Printf("Number of completed habits for %s: %d\n", name, result.Count)
	return result.Count, nil
}

func isCompletedInDay(name string, year int, month time.Month, day int) (bool, error) {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Define the filter to find the habit
	filter := bson.M{
		"name":                       name,
		"year":                       year,
		"month":                      month,
		fmt.Sprintf("boxes.%d", day): "■ ", // Filter for documents where the specified day's box contains the symbol "■"
	}

	// Perform the find operation to check if the document exists
	if err := Collection.FindOne(ctx, filter).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil // Document does not exist, habit is not completed on that day
		}
		return false, fmt.Errorf("error checking completion status for day %d: %v", day, err)
	}

	return true, nil // Document exists, habit is completed on that day
}
func countCompletedInDay(year int, month time.Month, day int) (int, error) {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Define the filter to find habits completed on the specified day
	filter := bson.M{
		"year":                       year,
		"month":                      month,
		fmt.Sprintf("boxes.%d", day): "■ ", // Filter for documents where the specified day's box contains the symbol "■"
	}

	// Perform the find operation to get all habits completed on the specified day
	cursor, err := Collection.Find(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("error finding habits completed on day %d: %v", day, err)
	}
	defer cursor.Close(ctx)

	// Count the number of habits completed on the specified day
	count := 0
	for cursor.Next(ctx) {
		count++
	}

	return count, nil
}
func getNumberOfHabitsInMonth(year int, month time.Month) (int, error) {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Define the filter to find the habits in the specified year and month
	filter := bson.M{
		"year":  year,
		"month": month,
	}

	// Perform the count operation to get the number of habits
	count, err := Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("error counting habits for year %d, month %s: %v", year, month, err)
	}

	return int(count), nil
}

func plotingGraph(completedData []int, year int, month time.Month) {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	data := completedData
	labels := make([]string, len(data))
	for i := range labels {
		labels[i] = fmt.Sprintf("%d", i+1)
	}

	// Determine the maximum value for y-axis
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

func main() {

	// Povezivanje sa bazom podataka
	ctx := context.TODO()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	mongoClient, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Printf("error connecting to MongoDB: %s\n", err.Error())
		return
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			fmt.Printf("error disconnecting from MongoDB: %s\n", err.Error())
		}
	}()

	t := time.Now()
	year := t.Year()
	month := t.Month()
	day := t.Day()
	numOfDays := daysInMonth(year, month)

	args := os.Args
	if len(args) >= 2 && args[1] == "add" {
		if len(args) == 3 {
			habitName := args[2]

			var boxes []string
			for i := 0; i < numOfDays; i++ {
				boxes = append(boxes, drawBox(1))
			}
			// Create a Habit object
			habit := Habit{
				Name:  habitName,
				Boxes: boxes, // Initialize array with 30 empty strings
				Year:  year,
				Month: month,
			}

			// Call addTask with the created Habit object
			addHabit(habit)
		} else {

			// Parse the month string into a time.Month value
			month, err := parseMonth(args[4])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Convert the integer argument to a string
			year, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("Error occurred:", err)
				return
			}
			numOfDays := daysInMonth(year, month)

			habitName := args[2]

			var boxes []string
			for i := 0; i < numOfDays; i++ {
				boxes = append(boxes, drawBox(1))
			}
			// Create a Habit object
			habit := Habit{
				Name:  habitName,
				Boxes: boxes, // Initialize array with 30 empty strings
				Year:  year,
				Month: month,
			}

			// Call addTask with the created Habit object
			addHabit(habit)
		}
	} else if len(args) >= 2 && args[1] == "all" {

		getAllHabits()
	} else if len(args) >= 2 && args[1] == "edit" {
		if len(args) == 4 {

			editHabitName(args[2], year, month, args[3])
		} else {
			// Parse the month string into a time.Month value
			month, err := parseMonth(args[4])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Convert the integer argument to a string
			id, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("Error occurred:", err)
				return
			}

			// Call the function
			editHabitName(args[2], id, month, args[5])
		}
	} else if len(args) >= 2 && args[1] == "delete" {
		if len(args) == 3 {

			deleteHabit(args[2], year, month)
		} else {
			// Parse the month string into a time.Month value
			month, err := parseMonth(args[4])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Convert the integer argument to a string
			year, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("Error occurred:", err)
				return
			}

			// Call the function
			deleteHabit(args[2], year, month)
		}
	} else if len(args) >= 2 && args[1] == "compleated" {
		if len(args) == 3 {

			compleatedHabit(args[2], year, month, day-1)
		} else {
			// Parse the month string into a time.Month value
			month, err := parseMonth(args[4])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Convert the integer argument to a string
			year, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("Error occurred:", err)
				return
			}

			day, err := strconv.Atoi(args[5])
			if err != nil {
				fmt.Println("Error occurred:", err)
				return
			}
			// Call the function
			compleatedHabit(args[2], year, month, day-1)
		}
	} else if len(args) >= 2 && args[1] == "info" {
		month, err := parseMonth(args[3])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Convert the integer argument to a string
		year, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Error occurred:", err)
			return
		}

		var compleatedData []int
		for i := 1; i <= numOfDays; i++ {
			completedCount, err := countCompletedInDay(year, month, i)
			if err != nil {
				fmt.Printf("Error counting completed habits for day %d: %v\n", i, err)
				continue
			}
			compleatedData = append(compleatedData, completedCount)
		}
		plotingGraph(compleatedData, year, month)
	} else {
		getCurrentHabits(year, month)
	}
}
