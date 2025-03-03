package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init(userDB *sql.DB) {
	db = userDB

	tables := []string{
		`CREATE TABLE IF NOT EXISTS image (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT,
			filepath TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS video (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT,
			filepath TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS gif (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT,
			filepath TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS model (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT,
			filepath TEXT
		)`,
	}

	for _, tableQuery := range tables {
		_, err := db.Exec(tableQuery)
		if err != nil {
			log.Fatal("Error creating tables", err)
		}
	}
}

func Insert(tableName, filename, filepath string) error {
	query := "INSERT INTO " + tableName + " (filename, filepath) VALUES (?, ?)"
	_, err := db.Exec(query, filename, filepath)
	return err
}

func GetPathByFilename(tableName, filename string) (string, error) {
	query := "SELECT filepath FROM " + tableName + " WHERE filename = ?"
	var filepath string
	err := db.QueryRow(query, filename).Scan(&filepath)
	return filepath, err
}

func DeleteRecordByFilename(tableName, filename string) error {
	query := "DELETE FROM " + tableName + " WHERE filename = ?"
	_, err := db.Exec(query, filename)
	return err
}

func ListAllFilenames(tableName string) ([]string, error) {
	query := "SELECT filename FROM " + tableName
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var filenames []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}
