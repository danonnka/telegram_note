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
	Run()
}

func Run() {
	IDFlag := &structFlag.StructMapCheck{
		IDPersonFlag: make(map[int64]*structFlag.BoolStruct),
	}

	logger := setupLogger("debug")
	cfg, err := config.MustLoad()
	if err != nil {
		logger.Error(err.Error())
	}
	//cfg.Token = "Write your token here"
	db, err := storage.NewSqliteStorage(cfg.StoragePath)
	if err != nil {
		logger.Error(err.Error())
	}
	token := cfg.Token
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		logger.Error("ошибка при создание бота" + err.Error())
	}

	bot.Debug = false

	u := tg.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		err := handler.MainHandler(bot, update, db, IDFlag)
		if err != nil {
			logger.Error(err.Error())
		}
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
