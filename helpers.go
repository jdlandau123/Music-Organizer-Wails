package main

import (
  "fmt"
  "path/filepath"
  "slices"
  "strings"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func PrintError(err error) {
  if err != nil {
    fmt.Println(err)
  }
}

func CheckFileExtension(filename string) bool {
  var supportedExtensions []string = []string{".mp3", ".flac", ".wav"}
  ext := filepath.Ext(filename)
  return slices.Contains(supportedExtensions, ext)
}

func GetFileFormat(filename string) string {
  ext := filepath.Ext(filename)
  ext = strings.Replace(ext, ".", "", 1)
  ext = strings.ToUpper(ext)
  return ext
}

func CheckAlbumExists(a Album) bool {
  db, err := sql.Open("sqlite3", dbPath)
  PrintError(err)
  
  rows, err := db.Query("SELECT COUNT(*) FROM albums WHERE Album = ? AND Artist = ?", a.Album, a.Artist)
  PrintError(err)

  defer db.Close()

  var count int
  for rows.Next() {
    err = rows.Scan(&count)
    PrintError(err)
  }
  return count > 0
}
