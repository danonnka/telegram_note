package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func CreateGorm(s string) (*Storage, error) {
	db, err := gorm.Open(sqlite.Open(s)) //gorm.Config вроде можно не писать
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Users{}, &Notes{})
	return &Storage{db: db}, nil
}

func (s Storage) AddUsers(ID int64) error {
	s.db.Create(&Users{ // Была ОШИБКА - не передал указать на структуру - ссылка нужна что бы gorm смог записать туда результат
		TelegramID: ID,
	})

	if s.db.Error != nil {
		return s.db.Error
	}
	return nil
}

func (s Storage) AddNote(id int64, note string) error {
	var us Users // нужно что бы создать обьект, в который загрузим данные из базы.
	err := s.db.First(&us, "Telegram_ID", id)
	if err != nil {
		return s.db.Error
	}
	//нужен us что бы взять от туда номер пользователя и добавитиь этот номер к заметки
	n := Notes{
		UserID: us.Id,
		Not:    note,
	}
	//возможно надо сохранить
	if err := s.db.Create(&n).Error; err != nil {
		return err
	}
	return nil
}

func (s Storage) ShowNote(id int64) ([]Notes, error) {
	var us Users
	if err := s.db.Preload("Notes").First(&us, "telegram_id = ?", id).Error; err != nil {
		return nil, err
	}
	return us.Note, nil
	//почему сдесь просто Note а в скобках []NoteS - S большая другая переменная
}

//как понять он загружает только Notes Б.Д. или сначало основную базу и через неё находит нужные notes
//preload = "Когда ты загрузишь пользователя, также заранее загрузи его связанные Notes".

//для удаления заметок
