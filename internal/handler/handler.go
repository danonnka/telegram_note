package handler

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegramNote/internal/button"
	"telegramNote/internal/storage"
	"telegramNote/internal/structFlag"
	"time"
)

func MainHandler(bot *tg.BotAPI, update tg.Update, db storage.Storage, flag *structFlag.StructMapCheck) error {

	if update.Message != nil {
		IDchat := update.Message.Chat.ID
		err := (db).AddUsers(IDchat)
		if err != nil {
			return err
		}

		if _, ok := flag.IDPersonFlag[IDchat]; !ok {
			flag.IDPersonFlag[IDchat] = &structFlag.BoolStruct{
				AddNoteFlag:   false,
				DeletNoteFlag: false,
			}
		}
		if update.Message.Text == "/start" {
			err := ShowButton(bot, update, button.RowButton, "выберите действие")
			if err != nil {
				return err
			}

		}
		if flag.IDPersonFlag[IDchat].AddNoteFlag == true {
			flag.IDPersonFlag[IDchat].AddNoteFlag = false
			text := update.Message.Text

			err := db.AddNote(IDchat, text)
			if err != nil {
				return err
			}

			message := tg.NewMessage(update.Message.Chat.ID, "ваша заметка создана")
			_, err = bot.Send(message)
			if err != nil {
				return err
			}
			err = ShowButton(bot, update, button.RowButton, "выберите действие")
			if err != nil {
				return err
			}
		}

		if flag.IDPersonFlag[IDchat].DeletNoteFlag == true {
			flag.IDPersonFlag[IDchat].DeletNoteFlag = false
			numberNote := update.Message.Text

			errWrongNumber := db.DeletNote(IDchat, numberNote)
			if errWrongNumber != nil {
				message := tg.NewMessage(update.Message.Chat.ID, errWrongNumber.Error())
				_, err := bot.Send(message)
				if err != nil {
					return err
				}
				err = ShowButton(bot, update, button.RowButton, "выберите действие")
				return nil
			}
			message := tg.NewMessage(update.Message.Chat.ID, "ваша заметка удалена")
			_, err := bot.Send(message)
			if err != nil {
				return err
			}
			err = ShowButton(bot, update, button.RowButton, "выберите действие")
			if err != nil {
				return err
			}
			return nil
		}

	}
	if update.CallbackQuery != nil {
		IDbutton := update.CallbackQuery.Message.Chat.ID

		if _, ok := flag.IDPersonFlag[IDbutton]; !ok {
			flag.IDPersonFlag[IDbutton] = &structFlag.BoolStruct{
				AddNoteFlag:   false,
				DeletNoteFlag: false,
			}
		}

		switch update.CallbackQuery.Data {
		case "createNote":
			flag.IDPersonFlag[IDbutton].DeletNoteFlag = false
			flag.IDPersonFlag[IDbutton].AddNoteFlag = true

			go func(IDbutton int64) {
				time.Sleep(time.Minute)
				if flag.IDPersonFlag[IDbutton].AddNoteFlag {
					flag.IDPersonFlag[IDbutton].AddNoteFlag = false
				}
			}(IDbutton)

			callback := tg.NewCallback(update.CallbackQuery.ID, "напишите вашу заметку и отправьте")
			_, _ = bot.Request(callback)

		case "showNote":
			flag.IDPersonFlag[IDbutton].AddNoteFlag = false
			flag.IDPersonFlag[IDbutton].DeletNoteFlag = false
			callback := tg.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				return err
			}
			resultNotes, err := db.ShowNote(IDbutton)
			if err != nil {
				return err
			}

			for _, v := range resultNotes {
				retur := fmt.Sprintf("%v) %v", v.ID, v.Not)
				sendBack := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, retur)
				_, err := bot.Send(sendBack)
				if err != nil {
					return err
				}
			}
			err = ShowButton(bot, update, button.RowButton, "выберите действие")
			if err != nil {
				return err
			}
		case "deleteNote":
			flag.IDPersonFlag[IDbutton].AddNoteFlag = false
			flag.IDPersonFlag[IDbutton].DeletNoteFlag = true

			go func(IDbutton int64) {
				time.Sleep(time.Minute)
				if flag.IDPersonFlag[IDbutton].DeletNoteFlag {
					flag.IDPersonFlag[IDbutton].DeletNoteFlag = false
				}
			}(IDbutton)

			callback := tg.NewCallback(update.CallbackQuery.ID, "напишите номер заметки для удаления")
			_, err := bot.Request(callback)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func ShowButton(bot *tg.BotAPI, update tg.Update, rowButton []tg.InlineKeyboardButton, text string) error {
	if update.Message != nil {
		message := tg.NewMessage(update.Message.Chat.ID, text)
		message.ReplyMarkup = tg.NewInlineKeyboardMarkup(rowButton)
		_, err := bot.Send(message)
		if err != nil {
			return err
		}
	} else if update.CallbackQuery != nil {
		message2 := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		message2.ReplyMarkup = tg.NewInlineKeyboardMarkup(rowButton)
		_, err := bot.Send(message2)
		if err != nil {
			return err
		}
	}
	return nil
}
