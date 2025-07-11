package main

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"myNote3/internal/config"
	"myNote3/internal/handler"
	"myNote3/internal/storage"
	"myNote3/internal/structFlag"
	"os"
)

func main() {
	token := "7632907437:AAFKWnffg7nxvu88xrQ4N_T8X_SBaAyEgbM"

	Run(token)
}

// возможно лучше перенести log и cfg в main
func Run(token string) {
	IDFlag := &structFlag.StructMapCheck{
		IDPersonFlag: make(map[int64]*structFlag.BoolStruct),
	} //интересно почему при записи второй структуры обычно адрес пакета structFlag и название
	//  так-же почему именно адрес на неё . Веть что бы изменятб вроде нужен &

	log := setupLogger("debug")                    //это запуск проверки err. Что бы писать log.___
	cfg := config.MustLoad()                       // это путь к sql
	db, err := storage.CreateGorm(cfg.StoragePath) //создали базу данных черее Gorm

	bot, err := tg.NewBotAPI(token)
	if err != nil {
		log.Error("ошибка при создание бота" + err.Error())
	}

	bot.Debug = false

	u := tg.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handler.MainHandler(bot, update, db, IDFlag) //интересно почему тут надо со звездой передовать db или ненадо
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "info":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}), //error and warn
		)
	case "debug":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}), //all
		)
	case "error":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}), //error
		)
	case "warn":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}), //error and warn
		)
	}
	return log
}
