package noteStruct

import (
	"database/sql"
	"log"
	"myNote3/config"
	"myNote3/internal/DB/getNotes"
)

type StructWithNote struct {
	CreateNoteMap map[int64]*Note
}
type Note struct {
	NoteForPerson    []string
	Check            bool
	FlagDeletMessage bool
	InsertedID       []int64
}

func (st *StructWithNote) InitFromDB() {
	checkOrCreateTable()
	data := getNotes.GetIDAndNoteSlice()

	for _, note := range data {
		if _, ok := st.CreateNoteMap[note.Id]; !ok {
			st.CreateNoteMap[note.Id] = &Note{}
		}
		st.CreateNoteMap[note.Id].NoteForPerson = append(st.CreateNoteMap[note.Id].NoteForPerson, note.Text)
	}
}

func checkOrCreateTable() {
	db, err := sql.Open("sqlite3", config.DB_PATH)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err, "can't close db connection")
		}
	}()

	sqlStmt := "CREATE TABLE IF NOT EXISTS notes (id INTEGER PRIMARY KEY AUTOINCREMENT,  telegram_ID INTEGER, note TEXT);"

	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("Table created successfully")

}
