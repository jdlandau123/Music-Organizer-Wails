package main

import (
	"context"
	"fmt"
  "os"
  "path/filepath"
  "errors"
  "strings"
  "encoding/json"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "github.com/wailsapp/wails/v2/pkg/runtime"
)

const dbPath = "./db.sqlite"

func printError(err error) {
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

// startup is called when the app starts. The context is saved so we can call the runtime methods
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
  printError(err)
  return dir
}

// config logic
type Config struct {
  Id int
  CollectionPath string
  DevicePath string
}

func (a *App) GetConfig() Config {
  db, err := sql.Open("sqlite3", dbPath)
  printError(err)

  rows, err := db.Query("SELECT * FROM config ORDER BY Id DESC LIMIT 1")

  defer rows.Close()

  var config Config

  for rows.Next() {
    var id int
    var collectionPath string
    var devicePath string

    err = rows.Scan(&id, &collectionPath, &devicePath)
    printError(err)
    config = Config{id, collectionPath, devicePath}
  }
  a.config = config
  return config
}

func (a *App) SetConfig(c Config) {
  db, err := sql.Open("sqlite3", dbPath)
  printError(err)

  _, err = db.Exec("INSERT INTO config VALUES (NULL, ?, ?)", c.CollectionPath, c.DevicePath)
  printError(err)

  db.Close()
}


// database logic
type Album struct {
  // Id int
  Album string
  Artist string
  FileFormat string
  Tracklist string
  IsOnDevice bool
}

type Track struct {
  Number int
  Title string
}

func (a *App) InitDb() {
  db, err := sql.Open("sqlite3", dbPath)
  printError(err)

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `albums` (Id INTEGER PRIMARY KEY, Album TEXT, Artist TEXT, FileFormat TEXT, Tracklist TEXT, IsOnDevice INTEGER)")
  printError(err)

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `config` (Id INTEGER PRIMARY KEY, CollectionPath TEXT, DevicePath TEXT)")
  printError(err)

  db.Close()
}

func (a *App) GetAlbums() []Album {
  db, err := sql.Open("sqlite3", dbPath)
  printError(err)
  
  rows, err := db.Query("SELECT * FROM albums ORDER BY Artist, Album")
  printError(err)
  
  var albums []Album
  for rows.Next() {
    var id int
    var album string
    var artist string
    var fileFormat string
    var tracklist string
    var isOnDevice bool

    err = rows.Scan(&id, &album, &artist, &fileFormat, &tracklist, &isOnDevice)
    printError(err)
    
    albums = append(albums, Album{album, artist, fileFormat, tracklist, isOnDevice})
  }
  return albums
}

func AddAlbumToDb(a Album) {
  db, err := sql.Open("sqlite3", dbPath)
  printError(err)
  
  rows, err := db.Query("SELECT COUNT(*) FROM albums WHERE Album = ? AND Artist = ?", a.Album, a.Artist)
  printError(err)

  var count int
  for rows.Next() {
    err = rows.Scan(&count)
    printError(err)
  }

  if count > 0 {
    _, err = db.Exec("UPDATE albums SET FileFormat = ?, Tracklist = ? WHERE Album = ? AND Artist = ?", a.FileFormat, a.Tracklist, a.Album, a.Artist)
    printError(err)
  } else {
    _, err = db.Exec("INSERT INTO albums VALUES (NULL, ?, ?, ?, ?, ?)", a.Album, a.Artist, a.FileFormat, a.Tracklist, a.IsOnDevice)
    printError(err)
  }
  defer db.Close()
}

func BuildTracklist(songs []os.DirEntry) string {
  var tracklist []Track
  for index, song := range songs {
    tracklist = append(tracklist, Track{index+1, song.Name()})
  }
  
  res, err := json.Marshal(tracklist)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(string(res))
  return string(res)
}

func (a *App) SyncMusicCollection() error {
  if _, err := os.Stat(a.config.CollectionPath); os.IsNotExist(err) {
    return errors.New("Music Collection Not Found")
  }

  artists, err := os.ReadDir(a.config.CollectionPath)
  printError(err)

  for _, artist := range artists {
    fmt.Println(artist.Name())
    artistDir := filepath.Join(a.config.CollectionPath, artist.Name())

    albums, err := os.ReadDir(artistDir)
    printError(err)

    for _, album := range albums {
      albumDir := filepath.Join(artistDir, album.Name())

      songs, err := os.ReadDir(albumDir)
      printError(err)

      ext := filepath.Ext(songs[0].Name())
      ext = strings.Replace(ext, ".", "", 1)
      ext = strings.ToUpper(ext)

      tracklist := BuildTracklist(songs)
      
      a := Album{album.Name(), artist.Name(), ext, tracklist, false}
      AddAlbumToDb(a)
    }
  }
  return nil
}


func (a *App) SyncDevice() error {
  if _, err := os.Stat(a.config.CollectionPath); os.IsNotExist(err) {
    return errors.New("Music Collection Not Found")
  }

  if _, err := os.Stat(a.config.DevicePath); os.IsNotExist(err) {
    return errors.New("Device Not Found")
  }
  return nil
}
