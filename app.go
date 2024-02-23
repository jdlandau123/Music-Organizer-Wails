package main

import (
	"context"
	"fmt"
  "os"
  "path/filepath"
  "errors"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

var supportedExtensions = []string{".mp3", ".flac", ".wav"}

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

func (a *App) InitDb() {
  db, err := sql.Open("sqlite3", dbPath)
  PrintError(err)

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `albums` (Id INTEGER PRIMARY KEY, Album TEXT, Artist TEXT, FileFormat TEXT, Tracklist TEXT, IsOnDevice INTEGER)")
  PrintError(err)

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `config` (Id INTEGER PRIMARY KEY, CollectionPath TEXT, DevicePath TEXT)")
  PrintError(err)

  defer db.Close()
}

func (a *App) SyncMusicCollection() error {
  if _, err := os.Stat(a.config.CollectionPath); os.IsNotExist(err) {
    return errors.New("Music Collection Not Found")
  }

  artists, err := os.ReadDir(a.config.CollectionPath)
  PrintError(err)

  for _, artist := range artists {
    fmt.Println(artist.Name())
    artistDir := filepath.Join(a.config.CollectionPath, artist.Name())

    albums, err := os.ReadDir(artistDir)
    PrintError(err)

    for _, album := range albums {
      albumDir := filepath.Join(artistDir, album.Name())

      songs, err := os.ReadDir(albumDir)
      PrintError(err)

      ext := GetFileFormat(songs[0].Name())

      tracklist := BuildTracklist(songs)
      
      a := Album{0, album.Name(), artist.Name(), ext, tracklist, false}
      AddAlbumToDb(a)
    }
  }
  return nil
}

func (a *App) ScanDevice() error {
  if _, err := os.Stat(a.config.DevicePath); os.IsNotExist(err) {
    return errors.New("Device Not Found")
  }
  artists, err := os.ReadDir(a.config.DevicePath)
  PrintError(err)
  
  for _, artist := range(artists) {
    artistDir := filepath.Join(a.config.DevicePath, artist.Name())
    albums, err := os.ReadDir(artistDir)
    PrintError(err)

    for _, album := range albums {
      a := GetAlbumByName(album.Name(), artist.Name())
      if a.Id > 0 {
        SetIsOnDevice(true, a)
      } else {
        a = Album{0, album.Name(), artist.Name(), "", "", true}
        AddAlbumToDb(a)
      }
    }
  }
  return nil
}

func (a *App) SyncDevice(ids []int) error {
  if _, err := os.Stat(a.config.CollectionPath); os.IsNotExist(err) {
    return errors.New("Music Collection Not Found")
  }

  if _, err := os.Stat(a.config.DevicePath); os.IsNotExist(err) {
    return errors.New("Device Not Found")
  }
 
  for _, id := range ids {
    album := GetAlbumById(id)
    a.TransferAlbum(album)
    SetIsOnDevice(true, album)
  }

  a.RemoveAlbumsFromDevice(ids)
  return nil
}
