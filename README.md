# Habit Tracker

The habit tracker gives us the ability to track our habits throughout the day and easily mark them as completed. This allows us to review our activity for the entire month or even multiple months, providing us with detailed insight into our productivity and engagement.

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
- `Habit completed name`: Mark the habit as completed for the current year and month.
- `Habit completed name 2024 March 20`: Mark the habit as completed for the specified year, month, and day.
- `Habit info 2024 March `: Retrieves a graph with all habits for this month and shows their completion throughout the days.

## Technologies Used

- Go
- MongoDB

## Installation

To run the project locally:

1. Clone this repository.
2. Install MongoDB and ensure it's running.
3. Install Go if not already installed.
4. Navigate to the project directory and run `go mod tidy && go build -o Habit`.
5. Start the Habit Tracker ./main


