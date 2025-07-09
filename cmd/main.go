package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"myNote3/internal/handler"
	"myNote3/internal/noteStruct"
)

func main() {
	token := "7632907437:AAFKWnffg7nxvu88xrQ4N_T8X_SBaAyEgbM"
	Run(token)
}
func Run(token string) {
	GlobalStruct := noteStruct.StructWithNote{
		map[int64]*noteStruct.Note{},
	}
	GlobalStruct.InitFromDB() // writing in GlobalStruct whole table

	bot, err := tg.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u := tg.NewUpdate(0)
	//u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	//for update := range updates {
	//		handler.MainHandler(bot, update, &GlobalStruct)

	//}
	for range 10 {
		update := <-updates
		handler.MainHandler(bot, update, &GlobalStruct)
	}
}

// сделать удаленеие (должно появиться перечень заметиок с номером) (ведите номер для удаления) (рефакториг)
