package main

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Config struct {
  Id int
  CollectionPath string
  DevicePath string
}

// opens file explorer and returns user selected dir path 
func (a *App) SelectDirectory() string {
  options := runtime.OpenDialogOptions{
    Title: "Select a directory",
  }
  dir, err := runtime.OpenDirectoryDialog(a.ctx, options)
  PrintError(err)
  return dir
}


func (a *App) GetConfig() Config {
  db, err := sql.Open("sqlite3", dbPath)
  PrintError(err)

  rows, err := db.Query("SELECT * FROM config ORDER BY Id DESC LIMIT 1")

  defer rows.Close()

  var config Config

  for rows.Next() {
    var id int
    var collectionPath string
    var devicePath string

    err = rows.Scan(&id, &collectionPath, &devicePath)
    PrintError(err)
    config = Config{id, collectionPath, devicePath}
  }
  a.config = config
  return config
}


func (a *App) SetConfig(c Config) {
  db, err := sql.Open("sqlite3", dbPath)
  PrintError(err)

  _, err = db.Exec("INSERT INTO config VALUES (NULL, ?, ?)", c.CollectionPath, c.DevicePath)
  PrintError(err)

  db.Close()
}

