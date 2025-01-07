package database

import (
	"database/sql"
	"fmt"
)

type ConnectionInfo struct {
	Name     string
	Port     int
	Host     string
	User     string
	Password string
	SSLMode  string
}

func NewPostgresConnection(info ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		info.Host, info.Port, info.User, info.Password, info.Name, info.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
