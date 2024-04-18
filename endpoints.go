package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

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
