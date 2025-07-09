package button

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var CreateButton = tg.NewInlineKeyboardButtonData("note", "createNote")
var ViewButton = tg.NewInlineKeyboardButtonData("show", "showNote")
var DeleteButton = tg.NewInlineKeyboardButtonData("delete", "deleteNote")
var RowButton = tg.NewInlineKeyboardRow(CreateButton, ViewButton, DeleteButton)
