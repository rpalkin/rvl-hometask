package db

import (
	"fmt"
	"time"
)

var (
	NotFoundError = fmt.Errorf("not found")
)

type DbInterface interface {
	GetBirthday(username string) (time.Time, error)
	PutBirthday(username string, value time.Time) error
}
