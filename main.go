package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
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
	// day := t.Day()
	numOfDays := daysInMonth(year, month)

	var boxes []string
	for i := 0; i < numOfDays; i++ {
		boxes = append(boxes, drawBox(1))
	}
	// simbol := "■ "
	// boxes[day] = simbol

	// for i := 0; i < len(boxes); i++ {
	// 	fmt.Printf("%s ", boxes[i])
	// }
	// fmt.Println()
	args := os.Args
	if len(args) >= 2 && args[1] == "add" {
		habitName := args[2]

		// Create a Habit object
		habit := Habit{
			Name:  habitName,
			Boxes: boxes, // Initialize array with 30 empty strings
			Year:  year,
			Month: month,
		}

		// Call addTask with the created Habit object
		addHabit(habit)
	} else if len(args) >= 2 && args[1] == "all" {

		getAllHabits()
	} else {

		getCurrentHabits(year, month)
	}

}
