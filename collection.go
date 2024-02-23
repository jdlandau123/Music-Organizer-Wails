package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	cp "github.com/otiai10/copy"
)

type Album struct {
	Id         int
	Album      string
	Artist     string
	FileFormat string
	Tracklist  string
	IsOnDevice bool
}

type Track struct {
	Number int
	Title  string
}

func (a *App) GetAlbums() []Album {
	db, err := sql.Open("sqlite3", dbPath)
	PrintError(err)

	rows, err := db.Query("SELECT * FROM albums ORDER BY Artist, Album")
	PrintError(err)

	var albums []Album
	for rows.Next() {
		var id int
		var album string
		var artist string
		var fileFormat string
		var tracklist string
		var isOnDevice bool

		err = rows.Scan(&id, &album, &artist, &fileFormat, &tracklist, &isOnDevice)
		PrintError(err)

		albums = append(albums, Album{id, album, artist, fileFormat, tracklist, isOnDevice})
	}
	return albums
}

func GetAlbumById(id int) Album {
	db, err := sql.Open("sqlite3", dbPath)
	PrintError(err)

	rows, err := db.Query("SELECT * FROM albums WHERE Id = ?", id)
	PrintError(err)

	defer db.Close()

	var album Album
	for rows.Next() {
		var id int
		var albumName string
		var artist string
		var fileFormat string
		var tracklist string
		var isOnDevice bool

		err = rows.Scan(&id, &albumName, &artist, &fileFormat, &tracklist, &isOnDevice)
		PrintError(err)

		album = Album{id, albumName, artist, fileFormat, tracklist, isOnDevice}
	}
	return album
}

func GetAlbumByName(albumName, artist string) Album {
	db, err := sql.Open("sqlite3", dbPath)
	PrintError(err)

	rows, err := db.Query("SELECT * FROM albums WHERE Album = ? AND Artist = ?", albumName, artist)
	PrintError(err)

	defer db.Close()

	var album Album
	for rows.Next() {
		var id int
		var albumName string
		var artist string
		var fileFormat string
		var tracklist string
		var isOnDevice bool

		err = rows.Scan(&id, &albumName, &artist, &fileFormat, &tracklist, &isOnDevice)
		PrintError(err)

		album = Album{id, albumName, artist, fileFormat, tracklist, isOnDevice}
	}
	return album
}

func AddAlbumToDb(a Album) {
	db, err := sql.Open("sqlite3", dbPath)
	PrintError(err)

	if CheckAlbumExists(a) {
		_, err = db.Exec("UPDATE albums SET FileFormat = ?, Tracklist = ? WHERE Album = ? AND Artist = ?", a.FileFormat, a.Tracklist, a.Album, a.Artist)
		PrintError(err)
	} else {
		_, err = db.Exec("INSERT INTO albums VALUES (NULL, ?, ?, ?, ?, ?)", a.Album, a.Artist, a.FileFormat, a.Tracklist, a.IsOnDevice)
		PrintError(err)
	}
	defer db.Close()
}

func BuildTracklist(songs []os.DirEntry) string {
	var tracklist []Track
	for index, song := range songs {
		if !CheckFileExtension(song.Name()) {
			tracklist = append(tracklist, Track{index + 1, song.Name()})
		}
	}

	res, err := json.Marshal(tracklist)
	PrintError(err)

	return string(res)
}

func SetIsOnDevice(val bool, album Album) {
	db, err := sql.Open("sqlite3", dbPath)
	PrintError(err)

	var dbVal int = 0
	if val {
		dbVal = 1
	}

	_, err = db.Exec("UPDATE albums SET IsOnDevice = ? WHERE Album = ? AND Artist = ?", dbVal, album.Album, album.Artist)
	PrintError(err)

	defer db.Close()
}

func (a *App) TransferAlbum(album Album) {
	srcDir := filepath.Join(a.config.CollectionPath, album.Artist, album.Album)
	destDir := filepath.Join(a.config.DevicePath, album.Artist, album.Album)

	os.MkdirAll(destDir, 0750)

	songs, err := os.ReadDir(srcDir)
	PrintError(err)

	for _, song := range songs {
		if _, err := os.Stat(filepath.Join(destDir, song.Name())); os.IsNotExist(err) {
			if CheckFileExtension(song.Name()) {
				songPath := filepath.Join(srcDir, song.Name())
				err := cp.Copy(songPath, filepath.Join(destDir, song.Name()))
				PrintError(err)
			}
		}
	}
}

func (a *App) RemoveAlbumsFromDevice(ids []int) {
	db, err := sql.Open("sqlite3", dbPath)
	PrintError(err)

	placeholders := strings.Trim(strings.Repeat("?,", len(ids)), ",")
	query := fmt.Sprintf("UPDATE albums SET IsOnDevice = 0 WHERE Id NOT IN (%s)", placeholders)

	stmt, err := db.Prepare(query)
	PrintError(err)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err = stmt.Exec(args...)
	PrintError(err)

	query = fmt.Sprintf("SELECT Album, Artist FROM albums WHERE Id NOT IN (%s)", placeholders)
	stmt, err = db.Prepare(query)
	PrintError(err)

	rows, err := stmt.Query(args...)
	PrintError(err)

	for rows.Next() {
		var album string
		var artist string

		err = rows.Scan(&album, &artist)
		PrintError(err)

		artistPath := filepath.Join(a.config.DevicePath, artist)
		albumPath := filepath.Join(artistPath, album)
		os.RemoveAll(albumPath)

		if dirs, _ := os.ReadDir(artistPath); len(dirs) == 0 {
			os.RemoveAll(artistPath)
		}
	}
}
