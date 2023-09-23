package service

import (
	"fmt"
	"github.com/pkg/errors"
	"hometask/internal/db"
	"time"
)

var (
	UserNotFoundError = fmt.Errorf("user not found")
)

type BirthdayService struct {
	Db db.DbInterface
}

func NewBirthdayService(db db.DbInterface) *BirthdayService {
	return &BirthdayService{
		Db: db,
	}
}

func (b *BirthdayService) PutBirthday(username string, birthday time.Time) error {
	return b.Db.PutBirthday(username, birthday)
}

func (b *BirthdayService) GetDaysToBirthday(username string) (int, error) {
	birthday, err := b.Db.GetBirthday(username)
	if err != nil {
		if errors.Is(err, db.NotFoundError) {
			return 0, UserNotFoundError
		}
		return 0, err
	}
	nextBirthday := time.Date(time.Now().Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.UTC)
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	if nextBirthday.Before(today) {
		nextBirthday = nextBirthday.AddDate(1, 0, 0)
	}
	diff := int(nextBirthday.Sub(today).Hours() / 24)
	return diff, nil
}
