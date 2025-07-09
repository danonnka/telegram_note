package handler

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"myNote3/internal/DB/saves"
	"myNote3/internal/button"
	"myNote3/internal/noteStruct"
	"strconv"
)

func MainHandler(bot *tg.BotAPI, update tg.Update, GlobalStruct *noteStruct.StructWithNote) {
	for _, vale := range GlobalStruct.CreateNoteMap {
		fmt.Println(vale)
	}
	if update.Message != nil {
		ID := update.Message.Chat.ID
		if _, ok := GlobalStruct.CreateNoteMap[ID]; !ok {
			GlobalStruct.CreateNoteMap[ID] = &noteStruct.Note{make([]string, 0), false, false, make([]int64, 0)}
		} else {
			if GlobalStruct.CreateNoteMap[ID].Check == true {
				text := update.Message.Text
				GlobalStruct.CreateNoteMap[ID].NoteForPerson = append(GlobalStruct.CreateNoteMap[ID].NoteForPerson, text)

				GlobalStruct.CreateNoteMap[ID].Check = false

				saves.NoteToDB(ID, text, GlobalStruct)
				//TODO добавить сообщение об успешном создании заметки
				message := tg.NewMessage(update.Message.Chat.ID, "ваша заметка создана")
				_, err := bot.Send(message)
				if err != nil {
					log.Println(err)
				}

			}
		}
		if update.Message.Text == "/start" {
			ShowButton(bot, update, button.RowButton, "выберите действие")
		}
		if GlobalStruct.CreateNoteMap[ID].FlagDeletMessage == true {
			test(bot, update, GlobalStruct, ID)
		}

	}
	if update.CallbackQuery != nil {
		ID := update.CallbackQuery.Message.Chat.ID
		if _, ok := GlobalStruct.CreateNoteMap[ID]; !ok {
			GlobalStruct.CreateNoteMap[ID] = &noteStruct.Note{make([]string, 0), false, false, make([]int64, 0)}
		}
		switch update.CallbackQuery.Data {
		case "createNote":
			GlobalStruct.CreateNoteMap[ID].Check = true
			callback := tg.NewCallback(update.CallbackQuery.ID, "напишите вашу заметку и отправьте")
			_, _ = bot.Request(callback)

		case "showNote":
			callback := tg.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println("Ошибка при отправке callback:", err)
			}

			if len(GlobalStruct.CreateNoteMap[ID].NoteForPerson) == 0 {
				message := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "вы еще не сохранили ни одной заметки")
				_, _ = bot.Request(message)
				return
			}

			for i, v := range GlobalStruct.CreateNoteMap[ID].NoteForPerson {
				returnSprinf2 := fmt.Sprintf("%v) %v", i, v)
				sendBack := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, returnSprinf2)
				_, err := bot.Send(sendBack)
				if err != nil {
					log.Panic(err, "Ошибка в callbackHandler")
				}
			}

		case "deleteNote":
			if len(GlobalStruct.CreateNoteMap[ID].NoteForPerson) == 0 {
				message := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, "отсутствуют заметки для удаления")
				_, _ = bot.Request(message)
				callback := tg.NewCallback(update.CallbackQuery.ID, "")
				if _, err := bot.Request(callback); err != nil {
					log.Println("Ошибка при отправке callback:", err)
				}
				return
			}
			GlobalStruct.CreateNoteMap[ID].FlagDeletMessage = true
			//другой флаг вырубить (ошибка если нажал две кнопки) GlobalStruct.CreateNoteMap[ID]. = false
			callback := tg.NewCallback(update.CallbackQuery.ID, "Напишите номер заметки, которую желаете удалить ")
			_, _ = bot.Request(callback)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

		}

	}
}

func ShowButton(bot *tg.BotAPI, update tg.Update, rowButton []tg.InlineKeyboardButton, text string) {
	if update.Message != nil {
		message := tg.NewMessage(update.Message.Chat.ID, text)
		message.ReplyMarkup = tg.NewInlineKeyboardMarkup(rowButton)
		_, err := bot.Send(message)
		if err != nil {
			log.Panic(err)
		}
	}
}

func test(bot *tg.BotAPI, update tg.Update, GlobalStruct *noteStruct.StructWithNote, ID int64) {
	text := update.Message.Text

	number, err := strconv.Atoi(text)
	if err != nil {
		message := tg.NewMessage(update.Message.Chat.ID, "Ошибка: В следующий раз введите номер заметки в виде числа")
		_, _ = bot.Send(message)
		return
	}

	note, ok := GlobalStruct.CreateNoteMap[ID]
	if !ok {
		//fmt.Println("Пользователь не найден")
		message := tg.NewMessage(update.Message.Chat.ID, "Пользователь не найден")
		_, err := bot.Send(message)
		if err != nil {
			log.Println(err)
		}
		return
	}

	if number < 0 || number >= len(note.NoteForPerson) {
		//fmt.Println("Неверный индекс")
		message := tg.NewMessage(update.Message.Chat.ID, "Неверный номер пользователя")
		_, err := bot.Send(message)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// Удаляем элемент по индексу из NoteForPerson
	note.NoteForPerson = append(note.NoteForPerson[:number], note.NoteForPerson[number+1:]...)
	note.FlagDeletMessage = false
	saves.DeleteNoteFromDB(GlobalStruct, ID, number)
	texts := fmt.Sprintf("Ваши заметка № %v удалена", number)
	message := tg.NewMessage(update.Message.Chat.ID, texts)
	_, err = bot.Send(message)
	if err != nil {
		log.Println(err)
	}
}
