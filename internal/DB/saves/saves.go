package saves

import (
	"database/sql"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"myNote3/internal/noteStruct"
	"strings"
)

func SaveMap(bot *tg.BotAPI, update tg.Update, globalStruct noteStruct.StructWithNote, ID int64) {

	//create()
	insert(globalStruct, ID)
	//load(globalStruct)
}

func create() {
	db, err := sql.Open("sqlite3", "./example2.db")

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err, "can't close db connection")
		}
	}()

	sqlStmt := "CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, idPerson INTEGER,  note TEXT)"

	_, err = db.Exec(sqlStmt)

	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("Table created successfully")

}

func insert(globalStruct noteStruct.StructWithNote, ID int64) {
	vv := globalStruct.CreateNoteMap[ID].NoteForPerson
	vvd := strings.Join(vv, "\n")
	db, err := sql.Open("sqlite3", "./example2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (idPerson, note) VALUES (?, ?)", ID, vvd)
	if err != nil {
		log.Fatalf("Error inserting row: %v", err)
	}

}

func load(globalStruct noteStruct.StructWithNote) {
	db, err := sql.Open("sqlite3", "./example2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
