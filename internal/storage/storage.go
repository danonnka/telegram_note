package storage

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
)

type Storage interface {
	AddUsers(ID int64) error
	AddNote(id int64, note string) error
	ShowNote(id int64) ([]Notes, error)
	DeletNote(id int64, numberNote string) error
}
type SqliteStorage struct {
	db *gorm.DB
}

func NewSqliteStorage(s string) (*SqliteStorage, error) {
	db, err := gorm.Open(sqlite.Open(s))
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Users{}, &Notes{})
	if err != nil {
		return nil, err
	}
	return &SqliteStorage{db: db}, nil
}

func (s SqliteStorage) AddUsers(ID int64) error {
	var user Users
	err := s.db.First(&user, "telegram_id = ?", ID).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	s.db.Create(&Users{
		TelegramID: ID,
	})

	if s.db.Error != nil {
		return s.db.Error
	}
	return nil
}

func (s SqliteStorage) AddNote(id int64, note string) error {
	var us Users
	tx := s.db.First(&us, "Telegram_ID", id)
	if tx.Error != nil {
		return tx.Error
	}

	n := Notes{
		UserID: us.Id,
		Not:    note,
	}

	if err := s.db.Create(&n).Error; err != nil {
		return err
	}
	return nil
}

func (s SqliteStorage) ShowNote(id int64) ([]Notes, error) {
	var us Users
	if err := s.db.Preload("ConnectNotes").First(&us, "telegram_id = ?", id).Error; err != nil {
		return nil, err
	}
	fmt.Print(us.ConnectNotes)
	return us.ConnectNotes, nil

}
func (s SqliteStorage) DeletNote(id int64, numberNote string) error {
	number, err := strconv.Atoi(numberNote)
	if err != nil {
		return fmt.Errorf("пожалуйста, введите номер заметки цифрами")
	}

	if number <= 0 {
		return fmt.Errorf("номер заметки должен быть положительным числом")
	}
	var user Users
	if err := s.db.First(&user, "telegram_id = ?", id).Error; err != nil {
		return err
	}
	result := s.db.Where("user_id = ? AND id = ?", user.Id, number).Delete(&Notes{})
	if result.RowsAffected == 0 {
		return fmt.Errorf("неверный номер заметки")
	}
	return nil
}
