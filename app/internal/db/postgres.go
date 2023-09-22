package db

import (
	"database/sql"
	"time"
)

type Postgres struct {
	dbConn *sql.DB
}

func (r *Postgres) GetBirthday(username string) (time.Time, error) {
	query := `SELECT birthday FROM rvl_birthdays WHERE username = $1`
	var birthday time.Time
	err := r.dbConn.QueryRow(query, username).Scan(&birthday)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, NotFoundError
		}
		return time.Time{}, err
	}
	return birthday, nil
}

func (r *Postgres) PutBirthday(username string, value time.Time) error {
	query := `INSERT INTO rvl_birthdays (username, birthday) VALUES ($1, $2) ON CONFLICT (username) DO UPDATE SET birthday = $2`
	_, err := r.dbConn.Exec(query, username, value)
	if err != nil {
		return err
	}
	return nil
}

func NewPostgres(dbConn *sql.DB) *Postgres {
	return &Postgres{
		dbConn: dbConn,
	}
}
