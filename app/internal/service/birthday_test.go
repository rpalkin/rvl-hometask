package service

import (
	"hometask/internal/db"
	"testing"
	"time"
)

type DbMock struct {
	Data map[string]time.Time
}

func (d *DbMock) GetBirthday(username string) (time.Time, error) {
	if val, ok := d.Data[username]; ok {
		return val, nil
	}
	return time.Time{}, db.NotFoundError
}

func (d *DbMock) PutBirthday(username string, value time.Time) error {
	d.Data[username] = value
	return nil
}

func TestBirthdayService_GetDaysToBirthday(t *testing.T) {
	type fields struct {
		Db db.DbInterface
	}
	type args struct {
		username string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "birthday is not today",
			fields:  fields{Db: &DbMock{Data: map[string]time.Time{"test": time.Now().AddDate(-5, 0, 150)}}},
			args:    args{username: "test"},
			want:    150,
			wantErr: false,
		},
		{
			name:    "birthday is today",
			fields:  fields{Db: &DbMock{Data: map[string]time.Time{"test": time.Now().AddDate(-10, 0, 0)}}},
			args:    args{username: "test"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "no user found",
			fields:  fields{Db: &DbMock{Data: map[string]time.Time{}}},
			args:    args{username: "test"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBirthdayService(tt.fields.Db)
			got, err := b.GetDaysToBirthday(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDaysToBirthday() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDaysToBirthday() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBirthdayService_PutBirthday(t *testing.T) {
	type fields struct {
		Db db.DbInterface
	}
	type args struct {
		username string
		birthday time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "birthday is yesterday",
			fields:  fields{Db: &DbMock{Data: map[string]time.Time{}}},
			args:    args{username: "test", birthday: time.Now().Add(-time.Hour * 24)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBirthdayService(tt.fields.Db)
			if err := b.PutBirthday(tt.args.username, tt.args.birthday); (err != nil) != tt.wantErr {
				t.Errorf("PutBirthday() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
