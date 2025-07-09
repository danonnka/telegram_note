package saves

import (
	"database/sql"
	"log"
	"myNote3/config"
	"myNote3/internal/noteStruct"
)

// кладу в базу данных
func NoteToDB(ID int64, text string, GlobalStruct *noteStruct.StructWithNote) {

	db, err := sql.Open("sqlite3", config.DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO notes (telegram_ID, note) VALUES (?, ?)", ID, text)
	if err != nil {
		log.Fatalf("  NoteToDB  Error inserting row: %v", err)
	}
	insertedID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("NoteToDB: error getting last insert ID: %v", err)
	}
	GlobalStruct.CreateNoteMap[ID].InsertedID = append(GlobalStruct.CreateNoteMap[ID].InsertedID, insertedID)
}
