package main

import (
	"context"
	"fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "github.com/wailsapp/wails/v2/pkg/runtime"
)

const dbPath = "./db.sqlite"

func handleError(err error) {
  if err != nil {
    fmt.Println(err)
  }
}

// App struct
type App struct {
	ctx context.Context
  config Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
  a.InitDb()
}

// opens file explorer and returns user selected dir path 
func (a *App) SelectDirectory() string {
  options := runtime.OpenDialogOptions{
    Title: "Select a directory",
  }
  dir, err := runtime.OpenDirectoryDialog(a.ctx, options)
  handleError(err)
  return dir
}

// database logic
type Config struct {
  Id int
  CollectionPath string
  DevicePath string
}

type Album struct {
  Id int
  Album string
  Artist string
  FileFormat string
  Tracklist []Track
  IsOnDevice bool
}

type Track struct {
  Number int
  Title string
}

func (a *App) InitDb() {
  db, err := sql.Open("sqlite3", dbPath)
  handleError(err)

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `albums` (Id INTEGER PRIMARY KEY, Album TEXT, Artist TEXT, FileFormat TEXT, Tracklist TEXT, IsOnDevice INTEGER)")
  handleError(err)

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `config` (Id INTEGER PRIMARY KEY, CollectionPath TEXT, DevicePath TEXT)")
  handleError(err)

  db.Close()
}

func (a *App) GetConfig() Config {
  db, err := sql.Open("sqlite3", dbPath)
  handleError(err)

  rows, err := db.Query("SELECT * FROM config ORDER BY Id DESC LIMIT 1")

  defer rows.Close()

  var config Config

  for rows.Next() {
    var id int
    var collectionPath string
    var devicePath string

    err = rows.Scan(&id, &collectionPath, &devicePath)
    handleError(err)
    config = Config{id, collectionPath, devicePath}
  }
  return config
}

func (a *App) SetConfig(c Config) {
  db, err := sql.Open("sqlite3", dbPath)
  handleError(err)

  _, err = db.Exec("INSERT INTO config VALUES (NULL, ?, ?)", c.CollectionPath, c.DevicePath)
  handleError(err)

  db.Close()
}


