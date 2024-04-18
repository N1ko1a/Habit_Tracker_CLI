package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
	"time"
)

// Kreiranje instance da bi smo se povezali sa bazom i imali interakciju
var mongoClient *mongo.Client

const url = "mongodb://localhost:27017"

func main() {

	// Povezivanje sa bazom podataka
	ctx := context.TODO()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)

	var err error
	mongoClient, err = mongo.Connect(ctx, opts)
	if err != nil {
		fmt.Printf("error connecting to MongoDB: %s\n", err.Error())
		return
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			fmt.Printf("error disconnecting from MongoDB: %s\n", err.Error())
		}
	}()

	//Trenutno vreme
	t := time.Now()
	year := t.Year()
	month := t.Month()
	day := t.Day()
	numOfDays := daysInMonth(year, month)

	args := os.Args //Omogucava koriscenje argumenata

	//Ako je argumenat prisutan
	if len(args) >= 2 {
		switch args[1] {
		case "add":
			if len(args) == 3 {
				habitName := args[2]

				var boxes []string // definisali smo niz kocki
				for i := 0; i < numOfDays; i++ {
					boxes = append(boxes, drawBox(1)) //dodajemo po jednu kocku za onoliko puta koliko imamo dana u mesecu
				}

				//Kreirali smo habit objekat
				habit := Habit{
					Name:  habitName,
					Boxes: boxes,
					Year:  year,
					Month: month,
				}

				//Dodajemo habit u bazu
				err := addHabit(habit)
				if err != nil { //ako dodje do greske printaj err koji nam funkcija addHabit daje
					fmt.Println("ErrorL ", err)
					return
				}
				getCurrentHabits(year, month)
			} else if len(args) == 5 {

				//Pripremamo podatke za objekat
				//Parsiramo mesec
				month, err := parseMonth(args[4])
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				//Persiramo godinu
				year, err := strconv.Atoi(args[3])
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}

				numOfDays := daysInMonth(year, month)
				habitName := args[2]

				//Kreiramo prazne kocke
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
				err = addHabit(habit)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				getCurrentHabits(year, month)
			} else {
				fmt.Println("Arguments error.")
			}
		case "all":
			if len(args) == 2 {
				getAllHabits()
			} else {
				fmt.Println("Arguments error.")
			}
		case "edit":
			if len(args) == 4 {
				err := editHabitName(args[2], year, month, args[3])
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				getCurrentHabits(year, month)
			} else if len(args) == 6 {
				//Parsiramo mesec
				month, err := parseMonth(args[4])
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				//Parsiramo godinu
				year, err := strconv.Atoi(args[3])
				if err != nil {
					fmt.Println("Error occurred:", err)
					return
				}
				// Call the function
				err = editHabitName(args[2], year, month, args[5])
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				getCurrentHabits(year, month)
			} else {
				fmt.Println("Arguments error.")
			}
		case "delete":
			if len(args) == 3 {
				err := deleteHabit(args[2], year, month)
				if err != nil {
					fmt.Println("Error", err)
					return
				}
				getCurrentHabits(year, month)
			} else if len(args) == 5 {
				//Parsiramo mesec
				month, err := parseMonth(args[4])
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				//Parsiramo godinu
				year, err := strconv.Atoi(args[3])
				if err != nil {
					fmt.Println("Error occurred:", err)
					return
				}

				err = deleteHabit(args[2], year, month)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				getCurrentHabits(year, month)
			} else {
				fmt.Println("Arguments error.")
			}
		case "completed":
			if len(args) == 3 {
				err := compleatedHabit(args[2], year, month, day)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				getCurrentHabits(year, month)
			} else if len(args) == 6 {
				//Parsiramo mesec
				month, err := parseMonth(args[4])
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				//Parsiramo godinu
				year, err := strconv.Atoi(args[3])
				if err != nil {
					fmt.Println("Error occurred:", err)
					return
				}

				//Parsiramo dan
				day, err := strconv.Atoi(args[5])
				if err != nil {
					fmt.Println("Error occurred:", err)
					return
				}

				err = compleatedHabit(args[2], year, month, day)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				getCurrentHabits(year, month)
			} else {
				fmt.Println("Arguments error.")

			}
		case "info":
			if len(args) == 4 {
				//Parsiramo mesec
				month, err := parseMonth(args[3])
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				//Parsiramo godinu
				year, err := strconv.Atoi(args[2])
				if err != nil {
					fmt.Println("Error occurred:", err)
					return
				}

				var compleatedData []int          // niz sa brojem ispunjenih navika po danima
				for i := 1; i <= numOfDays; i++ { // Prilazimo kroz dane i brojimo ispunjene navike pa dodajemo u niz
					completedCount, err := countCompletedInDay(year, month, i)
					if err != nil {
						fmt.Printf("Error counting completed habits for day %d: %v\n", i, err)
						continue
					}
					compleatedData = append(compleatedData, completedCount)
				}
				//Crtamo grafik
				plotingGraph(compleatedData, year, month)
			} else {
				fmt.Println("Arguments error.")
			}
		default:
			fmt.Println("Not valid command!")
		}
	} else {
		getCurrentHabits(year, month)
	}

}
