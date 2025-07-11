package storage

type Users struct {
	Id         int64   `gorm:"primaryKey"`
	TelegramID int64   `gorm:"uniqueIndex"`
	Note       []Notes `gorm:"foreignKey:UserID"`
} //думаю надо слайс заметок

type Notes struct {
	UserID int64  `gorm:"index;not null"` //нужен что бы установить связб между двумя таблицами
	Not    string `gorm:"not null"`
}

// foreignKey (внешний ключ) — это связь между двумя таблицами.
//Он указывает, какое поле в одной таблице связано с полем в другой таблице.
