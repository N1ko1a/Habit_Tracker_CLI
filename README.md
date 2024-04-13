# Habit Tracker

This is a simple habit tracker project that allows you to manage your habits through various commands. 

## Usage

### Commands:

- `Habit`: Show habits for the current year and month.
- `Habit all`: Show habits for all years and months.
- `Habit delete name`: Delete habit for the current year and month.
- `Habit delete name 2024 March`: Delete habit for the specified year and month.
- `Habit edit name newName`: Change habit name to new name for the current year and month.
- `Habit edit name 2024 March newName`: Change habit name to new name for the specified year and month.
- `Habit add name 2024 March`: Add habit for the specified year and month.
- `Habit add name`: Add habit for the current year and month.
- `Habit compleated name`: Mark the habit as completed for the current year and month.
- `Habit compleated name 2024 March 20`: Mark the habit as completed for the specified year, month, and day.

## Technologies Used

- Go
- MongoDB

## Installation

To run the project locally:

1. Clone this repository.
2. Install MongoDB and ensure it's running.
3. Install Go if not already installed.
4. Navigate to the project directory and run `go build main.go`.
5. Start the Habit Tracker ./main


