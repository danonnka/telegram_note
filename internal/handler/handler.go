package handler

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"myNote3/internal/button"
	"myNote3/internal/storage"
	"myNote3/internal/structFlag"
)

// идея для будущего. Возможно заменить флажки на поиск по sql
func MainHandler(bot *tg.BotAPI, update tg.Update, db *storage.Storage, flag *structFlag.StructMapCheck) {
	//как заменить флажок - не получилось. Пришлось заново создавать структуру
	if update.Message != nil {
		IDchat := update.Message.Chat.ID
		err := db.AddUsers(IDchat) //метод который добавляет пользователей по id
		_ = err                    //отправить ошибку потом в main

		if _, ok := flag.IDPersonFlag[IDchat]; !ok {
			flag.IDPersonFlag[IDchat] = &structFlag.BoolStruct{ //как понять что второе щзначение где флаги это полное имя типа структуры
				CheckFlag: false,
				DeletFlag: false,
			}
		} else {
			if flag.IDPersonFlag[IDchat].CheckFlag == true {
				//ошибка если у нас ещё нет никаких записей - а потом мы менаем полня на ture у ничего (nil)
				text := update.Message.Text
				flag.IDPersonFlag[IDchat].CheckFlag = false

				errs := db.AddNote(IDchat, text)
				_ = errs //вернуть ошибку в main

				message := tg.NewMessage(update.Message.Chat.ID, "ваша заметка создана")
				_, err := bot.Send(message)
				if err != nil {
					log.Println(err)
				}

			}

			if update.Message.Text == "/start" {
				ShowButton(bot, update, button.RowButton, "выберите действие")
			}
		}
		//if flag.IDPersonFlag[IDchat].DeletFlag == true {
		//test(bot, update, GlobalStruct, IDchat) //засунуть туда метод для удаления из базы данных
		//}
	}
	if update.CallbackQuery != nil {
		IDbutton := update.CallbackQuery.Message.Chat.ID

		switch update.CallbackQuery.Data {
		case "createNote":
			flag.IDPersonFlag[IDbutton].CheckFlag = true
			callback := tg.NewCallback(update.CallbackQuery.ID, "напишите вашу заметку и отправьте")
			_, _ = bot.Request(callback)

		case "showNote": //сюда метод показа заметок из базы
			callback := tg.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				log.Println("Ошибка при отправке callback:", err)
			}
			resultNotes, err := db.ShowNote(IDbutton)
			_ = err

			for i, v := range resultNotes {
				retur := fmt.Sprintf("%v) %v", i, v)
				sendBack := tg.NewMessage(update.CallbackQuery.Message.Chat.ID, retur)
				_, err := bot.Send(sendBack)
				if err != nil {
					log.Panic(err, "Ошибка в callbackHandler")
				}
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
