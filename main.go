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

// Kreiranje instance da bi smo se povezali sa bazom i imali interakciju
var mongoClient *mongo.Client

const url = "mongodb://localhost:27017"

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

// ENDPOITNS
// Dodavanje navike
func addHabit(habit Habit) error {
	ctx := context.TODO()                            //Vraca parzan context koji sluzi kao placeholder
	Database := mongoClient.Database("HabitTracker") //Baza koju koristimo ili kreiramo ako ne postoji
	Collection := Database.Collection("Habits")      //Kolekcija koju koristimo

	//Proveravamo da li postoji navika sa ovim imenom,godinom i mesecom
	filter := bson.D{{Key: "name", Value: habit.Name}, {Key: "year", Value: habit.Year}, {Key: "month", Value: habit.Month}}
	existingDoc := Collection.FindOne(ctx, filter)

	if existingDoc.Err() != nil && existingDoc.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("error while checking for existing documents: %v", existingDoc.Err())
	}

	if existingDoc.Err() == nil {
		fmt.Printf("Skipping insertion for habit - Name: %s as it already exists in the %d %d.\n", habit.Name, habit.Year, habit.Month)
		return nil
	} else {
		_, err := Collection.InsertOne(ctx, bson.D{
			{Key: "name", Value: habit.Name},
			{Key: "boxes", Value: habit.Boxes},
			{Key: "year", Value: habit.Year},
			{Key: "month", Value: habit.Month},
		})

		return err
	}
}

// Preuzimanje svih navika
func getAllHabits() {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	sortOptions := bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}

	//Koristimo Find sa bson.D{} sto znaci da ce da preuzme sve dokumente iz baze i te dokumente smo sortirali
	//cursor je rezultat upita i eventualnu gresku err
	cursor, err := Collection.Find(ctx, bson.D{}, options.Find().SetSort(sortOptions))
	if err != nil {
		fmt.Println("Error querying from the database:", err)
		return
	}
	//Obavezno zatvaramo cursor
	defer cursor.Close(ctx)

	var habits []Habit             //deklarisemo promenljivu habits koja je rezervisana za skladistenje niza Habit structure
	err = cursor.All(ctx, &habits) //Preuzima sve rezultate iz kursora i dekodira ih u navedenu strukturu podataka koja je u ovom slucaju niz habits. Prosledjujemo referencu niza habits tako da metoda moze direktno upisati rezultat u niz
	if err != nil {
		fmt.Println("Error decoding habits from the cursor:", err)
	}

	//Printovanje podataka
	var prevMonth time.Month
	var prevYear int

	for _, habit := range habits {
		//Proveravamo da li se razlikuje godina i mesec da bih razdvojili navike po njima
		if habit.Month != prevMonth || habit.Year != prevYear {
			fmt.Println()
			fmt.Println()
			fmt.Printf("%s %d\n", habit.Month, habit.Year)
			fmt.Println()
		}

		nameWidth := 15 //Duzina polja za ime navike
		//%-pocetak formatne specifikacije, - levo poravnjivanje teksta * oznacava da ce sirina polja biti odredjena s promenljiva
		fmt.Printf("%-*s", nameWidth, habit.Name)

		//Printovanje kockica za odredjenu naviku
		for i := 0; i < len(habit.Boxes); i++ {
			fmt.Printf("%s ", habit.Boxes[i])
		}
		fmt.Println()

		prevMonth = habit.Month
		prevYear = habit.Year
	}
	fmt.Println()
	fmt.Println()
}

// Preuzimanje trenutnih navika
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

	//Dekodiranje rezultata
	var habits []Habit
	err = cursor.All(ctx, &habits)
	if err != nil {
		fmt.Println("Error decoding habits from the cursor:", err)
		return
	}

	//Printovanja navika
	fmt.Println()
	fmt.Println()
	fmt.Printf("%s %d\n", month, year)
	fmt.Println()

	for _, habit := range habits {
		nameWidth := 15
		fmt.Printf("%-*s", nameWidth, habit.Name)
		for i := 0; i < len(habit.Boxes); i++ {
			fmt.Printf("%s ", habit.Boxes[i])
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
}

// Izmena imena navike
func editHabitName(name string, year int, month time.Month, newName string) error {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// filter := bson.M{}
	// filter["name"] = name
	// filter["year"] = year
	// filter["month"] = month

	//Drugi nacin za pisanja filtera
	filter := bson.D{{Key: "name", Value: name}, {Key: "year", Value: year}, {Key: "month", Value: month}}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: newName}}}}
	//Ako koristimo bson.M
	//   update := bson.M{
	//     "$set": bson.M{"name": newName},
	// }

	// _ sa ovim ignorisemo povratnu vrednost
	_, err := Collection.UpdateOne(ctx, filter, update) // Update jedan element koji nalazimo uz pomoc filtera i primenjujemo update na njega
	if err != nil {
		return fmt.Errorf("error updating habit: %v", err)
	}
	return nil
}

// Ispunjena navika
func compleatedHabit(name string, year int, month time.Month, day int) error {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	//Proveravamo da li je dan validan
	if !isValidDay(year, month, day) {
		return fmt.Errorf("Invalid day: %d for year: %d and month: %s", day, year, month)
	}

	//Proveravamo da li postoji habit sa ovim imenom
	filter := bson.D{{Key: "name", Value: name}, {Key: "year", Value: year}, {Key: "month", Value: month}}
	habitExists := Collection.FindOne(ctx, filter)

	if habitExists.Err() != nil && habitExists.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("Error while checking for existing documents: %v", habitExists.Err())
	}
	if habitExists.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("Habit with name %q does not exist", name)
	}

	symbol := "■ "
	update := bson.D{{Key: "$set", Value: bson.D{{Key: fmt.Sprintf("boxes.%d", day-1), Value: symbol}}}}

	_, err := Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Error compleating habit: %v ", err)
	}

	return nil

}

// Brisanje navike
func deleteHabit(name string, year int, month time.Month) error {
	ctx := context.TODO()
	Database := mongoClient.Database("HabitTracker")
	Collection := Database.Collection("Habits")

	// Check if a habit with the given name exists
	habitFilter := bson.D{{Key: "name", Value: name}, {Key: "year", Value: year}, {Key: "month", Value: month}}
	habitExists := Collection.FindOne(ctx, habitFilter)

	if habitExists.Err() != nil && habitExists.Err() != mongo.ErrNoDocuments {
		return fmt.Errorf("Error while checking for existing documents: %v", habitExists.Err())
	}
	if habitExists.Err() == mongo.ErrNoDocuments {
		return fmt.Errorf("Habit with name %q does not exist", name)
	}

	// Perform the deletion operation
	_, err := Collection.DeleteOne(ctx, habitFilter)
	if err != nil {
		return fmt.Errorf("error deleting habit: %v", err)
	}

	return nil
}

//GRAFIK

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
