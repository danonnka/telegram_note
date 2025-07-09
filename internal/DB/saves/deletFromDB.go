package saves

import (
	"database/sql"
	"log"
	"myNote3/config"
	"myNote3/internal/noteStruct"
)

func DeleteNoteFromDB(globalStruct *noteStruct.StructWithNote, ID int64, number int) {
	// Получаем ID заметки для удаления
	noteIDs := globalStruct.CreateNoteMap[ID].InsertedID
	if len(noteIDs) <= number {
		log.Printf("Неверный индекс заметки")
		return
	}

	// ID заметки, которую нужно удалить
	noteID := noteIDs[number]

	// Открываем базу данных
	db, err := sql.Open("sqlite3", config.DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Выполняем удаление
	_, err = db.Exec("DELETE FROM notes WHERE id = ?", noteID)
	if err != nil {
		log.Fatalf("Error deleting note with ID %d: %v", noteID, err)
	}
	log.Printf("Note with ID %d deleted", noteID)

	// Удаляем ID из слайса в структуре
	globalStruct.CreateNoteMap[ID].InsertedID = append(globalStruct.CreateNoteMap[ID].InsertedID[:number], globalStruct.CreateNoteMap[ID].InsertedID[number+1:]...)
}
