package service

import "time"

type BirthdayServiceInterface interface {
	PutBirthday(username string, birthday time.Time) error
	GetDaysToBirthday(username string) (int, error)
}
