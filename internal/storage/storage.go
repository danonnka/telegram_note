package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage interface {
	AddUsers(ID int64) error
	AddNote(id int64, note string) error
	ShowNote(id int64) ([]Notes, error)
}
type SqliteStorage struct {
	db *gorm.DB //теперь если мы изменим на другую таблицу это будет всё равно работать
} //              так-как эта структура реализует все необходимые методы для интерфейса

func NewSqliteStorage(s string) (*SqliteStorage, error) {
	db, err := gorm.Open(sqlite.Open(s))
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Users{}, &Notes{})
	return &SqliteStorage{db: db}, nil
}

func (s SqliteStorage) AddUsers(ID int64) error {
	var user Users
	if err := s.db.First(&user, "telegram_id = ?", ID).Error; err == nil {
		return nil
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
	var us Users // нужно что бы создать обьект, в который загрузим данные из базы.
	tx := s.db.First(&us, "Telegram_ID", id)
	if tx.Error != nil {
		return s.db.Error
	}
	//нужен us что бы взять от туда номер пользователя и добавитиь этот номер к заметки
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

//preload = "Когда ты загрузишь пользователя, также заранее загрузи его связанные Notes".

//для удаления заметок
