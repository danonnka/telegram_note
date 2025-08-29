package main

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"telegramNote/internal/config"
	"telegramNote/internal/handler"
	"telegramNote/internal/storage"
	"telegramNote/internal/structFlag"
)

func main() {
	cfg, err := config.MustLoad()
	logger := setupLogger("debug")

	err = Run(cfg)

	if err != nil {
		logger.Error(err.Error())
	}
}

func Run(cfg *config.Config) error {
	IDFlag := &structFlag.StructMapCheck{
		IDPersonFlag: make(map[int64]*structFlag.BoolStruct),
	}

	//cfg.Token = "Write your token here"

	db, err := storage.NewSqliteStorage(cfg.StoragePath)
	if err != nil {
		return err
	}
	token := cfg.Token
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("ошибка при создании бота: %w", err)
	}

	bot.Debug = false

	u := tg.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		err := handler.MainHandler(bot, update, db, IDFlag)
		if err != nil {
			return err
		}
	}
	return nil
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
