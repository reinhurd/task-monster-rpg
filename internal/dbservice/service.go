package dbservice

import (
	"database/sql"
	"fmt"
)

type DBService struct {
	db *sql.DB
}

func NewDBService(connStr string) (*DBService, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &DBService{db: db}, nil
}

func (s *DBService) Save(data interface{}) error {
	// TODO: Implement save method
	return nil
}

func (s *DBService) Select(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	return rows, nil
}

func (s *DBService) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
