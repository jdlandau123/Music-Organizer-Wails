package main

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	_ "modernc.org/sqlite"
)

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

func CheckAlbumExists(a Album) (bool, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return false, errors.New("Error connecting to database")
	}

	rows, err := db.Query("SELECT COUNT(*) FROM albums WHERE Album = ? AND Artist = ?", a.Album, a.Artist)
	if err != nil {
		fmt.Println(err)
		return false, errors.New("Error querying database")
	}

	defer db.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			fmt.Println(err)
		}
	}
	return count > 0, nil
}
