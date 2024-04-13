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
	Name  string   `json:"name"`
	Boxes []string `json:"boxes"`
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

func addTask(habit Habit) error {
	ctx := context.TODO()                            //vraca prazan context koji sluzi kao placeholder
	Database := mongoClient.Database("HabitTracker") //Ovo je baza, koristimo instancu inzad koju smo kreirali
	Collection := Database.Collection("Habits")      //Ovo je colekcija, koristimo umesto instance ime baze

	_, err := Collection.InsertOne(ctx, bson.D{
		{Key: "name", Value: habit.Name},
		{Key: "boxes", Value: habit.Boxes},
	})

	return err

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
	args := os.Args
	if len(args) >= 2 {
		habitName := args[2]

		// Create a Habit object
		habit := Habit{
			Name:  habitName,
			Boxes: boxes, // Initialize array with 30 empty strings
		}

		// Call addTask with the created Habit object
		addTask(habit)
	}

}
