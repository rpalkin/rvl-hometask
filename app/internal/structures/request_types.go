package structures

import (
	"fmt"
	"hometask/pkg/json_utils"
	"time"
)

type PutBirthdayRequest struct {
	DateOfBirth json_utils.JSONTime `json:"dateOfBirth"`
}

func (r *PutBirthdayRequest) Validate() error {
	if r.DateOfBirth.After(time.Now().Truncate(time.Hour * 24).Add(-time.Second)) {
		return fmt.Errorf("birthday must be before today")
	}
	return nil
}
