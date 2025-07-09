package getNotes

import (
	"database/sql"
	"log"
	"myNote3/config"
)

type SelectNote struct {
	Id   int64
	Text string
}

func GetIDAndNoteSlice() []SelectNote {
	SelectSlice := []SelectNote{}
	// Open the database connection
	db, err := sql.Open("sqlite3", config.DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query to select all users
	rows, err := db.Query("SELECT telegram_ID, note FROM notes ")
	if err != nil {
		log.Fatalf("Error querying rows: %v", err)
	}
	defer rows.Close()

	// Iterate through the rows
	for rows.Next() {
		var idper int64 // в эту переменную запист телеграм Айди err = rows.Scan(&idper, &superNote)
		var superNote string

		// Scan the row into variables
		err = rows.Scan(&idper, &superNote)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		SelectSlice = append(SelectSlice, SelectNote{idper, superNote})

		// Print the retrieved data

	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		log.Fatalf("Error during row iteration: %v", err)
	}
	return SelectSlice
}
