package storage

type Users struct {
	ID           int64   `gorm:"primaryKey"`
	TelegramID   int64   `gorm:"uniqueIndex"`
	ConnectNotes []Notes `gorm:"foreignKey:UserID"`
}

type Notes struct {
	ID     int64  `gorm:"primaryKey"`
	UserID int64  `gorm:"index;not null"`
	Not    string `gorm:"not null"`
}
