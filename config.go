package main

import (
	"database/sql"
	"errors"
	"fmt"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

type Config struct {
	Id             int
	CollectionPath string
	DevicePath     string
}

// opens file explorer and returns user selected dir path
func (a *App) SelectDirectory() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Select a directory",
	}
	dir, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Can't open directory")
	}
	return dir, nil
}

func (a *App) GetConfig() (Config, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return Config{}, errors.New("Error connecting to database")
	}

	rows, err := db.Query("SELECT * FROM config ORDER BY Id DESC LIMIT 1")

	defer rows.Close()

	var config Config

	for rows.Next() {
		var id int
		var collectionPath string
		var devicePath string

		err = rows.Scan(&id, &collectionPath, &devicePath)
		if err != nil {
			fmt.Println(err)
		}
		config = Config{id, collectionPath, devicePath}
	}
	a.config = config
	return config, nil
}

func (a *App) SetConfig(c Config) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error connecting to database")
	}

	_, err = db.Exec("INSERT INTO config VALUES (NULL, ?, ?)", c.CollectionPath, c.DevicePath)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error inserting row into config table")
	}

	defer db.Close()
	return nil
}
