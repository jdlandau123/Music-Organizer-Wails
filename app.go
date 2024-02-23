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
  if err != nil {
    fmt.Println("Error initializing database connection")
  }

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `albums` (Id INTEGER PRIMARY KEY, Album TEXT, Artist TEXT, FileFormat TEXT, Tracklist TEXT, IsOnDevice INTEGER)")
  if err != nil {
    fmt.Println("Error creating albums table")
  }

  _, err = db.Exec("CREATE TABLE IF NOT EXISTS `config` (Id INTEGER PRIMARY KEY, CollectionPath TEXT, DevicePath TEXT)")
  if err != nil {
    fmt.Println("Error creating config table")
  }

  defer db.Close()
}

func (a *App) SyncMusicCollection() error {
  if _, err := os.Stat(a.config.CollectionPath); os.IsNotExist(err) {
    return errors.New("Music Collection Not Found")
  }

  artists, err := os.ReadDir(a.config.CollectionPath)
  if err != nil {
    return errors.New("Error reading music collection directory")
  }

  for _, artist := range artists {
    fmt.Println(artist.Name())
    artistDir := filepath.Join(a.config.CollectionPath, artist.Name())

    albums, err := os.ReadDir(artistDir)
    if err != nil {
      fmt.Printf("Error reading directory for: %s", artist.Name())
      continue
    }

    for _, album := range albums {
      albumDir := filepath.Join(artistDir, album.Name())

      songs, err := os.ReadDir(albumDir)
      if err != nil {
        fmt.Printf("Error reading directory for: %s - %s", artist.Name(), album.Name())
        continue
      }

      ext := GetFileFormat(songs[0].Name())

      tracklist, err := BuildTracklist(songs)
      if err != nil {
        fmt.Printf("Building tracklist was unsuccessful for %s", album.Name())
      }
      
      a := Album{0, album.Name(), artist.Name(), ext, tracklist, false}
      _, err = AddAlbumToDb(a)
      if err != nil {
        fmt.Printf("Error adding album %s to database: %s", album.Name(), err)
        continue
      }
    }
  }
  return nil
}

func (a *App) ScanDevice() error {
  if _, err := os.Stat(a.config.DevicePath); os.IsNotExist(err) {
    return errors.New("Device Not Found")
  }

  artists, err := os.ReadDir(a.config.DevicePath)
  if err != nil {
    return errors.New("Error reading device directory")
  }
  
  var onDeviceIds []int

  for _, artist := range(artists) {
    artistDir := filepath.Join(a.config.DevicePath, artist.Name())
    albums, err := os.ReadDir(artistDir)
    if err != nil {
      fmt.Printf("Error reading directory for: %s", artist.Name())
      continue
    }

    for _, album := range albums {
      a, err := GetAlbumByName(album.Name(), artist.Name())
      if err != nil {
        fmt.Printf("No album found - %s", album.Name())
      }
      
      if a.Id > 0 {
        SetIsOnDevice(true, a)
        onDeviceIds = append(onDeviceIds, a.Id)
      } else {
        a = Album{0, album.Name(), artist.Name(), "", "", true}
        newId, err := AddAlbumToDb(a)
        if err != nil {
          fmt.Printf("Error adding album %s to database: %s", album.Name(), err)
          continue
        }
        onDeviceIds = append(onDeviceIds, newId)
      }
    }
  }

  a.RemoveAlbumsFromDevice(onDeviceIds)
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
    album, err := GetAlbumById(id)
    if err != nil {
      fmt.Printf("Failed to fetch album with Id = %v", id)
    }

    a.TransferAlbum(album)
    SetIsOnDevice(true, album)
  }

  a.RemoveAlbumsFromDevice(ids)
  return nil
}
