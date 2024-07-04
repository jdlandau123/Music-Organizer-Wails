package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
	_ "modernc.org/sqlite"
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

func (a *App) GetAlbums() ([]Album, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return []Album{}, errors.New("Error connecting to database")
	}

	rows, err := db.Query("SELECT * FROM albums ORDER BY Artist, Album")
	if err != nil {
		fmt.Println(err)
		return []Album{}, errors.New("Error querying database")
	}

	var albums []Album
	for rows.Next() {
		var id int
		var album string
		var artist string
		var fileFormat string
		var tracklist string
		var isOnDevice bool

		err = rows.Scan(&id, &album, &artist, &fileFormat, &tracklist, &isOnDevice)
		if err != nil {
			fmt.Println(err)
		}

		albums = append(albums, Album{id, album, artist, fileFormat, tracklist, isOnDevice})
	}
	return albums, nil
}

func GetAlbumById(id int) (Album, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return Album{}, errors.New("Error connecting to database")
	}

	rows, err := db.Query("SELECT * FROM albums WHERE Id = ?", id)
	if err != nil {
		fmt.Println(err)
		return Album{}, errors.New("Error querying database")
	}

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
		if err != nil {
			fmt.Println(err)
		}

		album = Album{id, albumName, artist, fileFormat, tracklist, isOnDevice}
	}
	return album, nil
}

func GetAlbumByName(albumName, artist string) (Album, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return Album{}, errors.New("Error connecting to database")
	}

	rows, err := db.Query("SELECT * FROM albums WHERE Album = ? AND Artist = ?", albumName, artist)
	if err != nil {
		fmt.Println(err)
		return Album{}, errors.New("Error querying database")
	}

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
		if err != nil {
			fmt.Println(err)
		}

		album = Album{id, albumName, artist, fileFormat, tracklist, isOnDevice}
	}
	return album, nil
}

func AddAlbumToDb(a Album) (int, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("Error connecting to database")
	}

	exists, err := CheckAlbumExists(a)
	if err != nil {
		fmt.Println(err)
	}

	var id int

	if exists {
		err := db.QueryRow(
			"UPDATE albums SET FileFormat = ?, Tracklist = ? WHERE Album = ? AND Artist = ? RETURNING Id",
			a.FileFormat,
			a.Tracklist,
			a.Album,
			a.Artist,
		).Scan(&id)
		if err != nil {
			fmt.Println(err)
			return id, errors.New("Error updating album")
		}
	} else {
		err := db.QueryRow(
			"INSERT INTO albums VALUES (NULL, ?, ?, ?, ?, ?) RETURNING Id",
			a.Album,
			a.Artist,
			a.FileFormat,
			a.Tracklist,
			a.IsOnDevice,
		).Scan(&id)
		if err != nil {
			fmt.Println(err)
			return id, errors.New("Error adding album")
		}
	}
	defer db.Close()
	return id, nil
}

func BuildTracklist(songs []os.DirEntry) (string, error) {
	var tracklist []Track
	for index, song := range songs {
		if CheckFileExtension(song.Name()) {
			tracklist = append(tracklist, Track{index + 1, song.Name()})
		}
	}

	res, err := json.Marshal(tracklist)
	if err != nil {
		return "", errors.New("Error marshalling json")
	}

	return string(res), nil
}

func SetIsOnDevice(val bool, album Album) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error connecting to database")
	}

	var dbVal int = 0
	if val {
		dbVal = 1
	}

	_, err = db.Exec(
		"UPDATE albums SET IsOnDevice = ? WHERE Album = ? AND Artist = ?",
		dbVal,
		album.Album,
		album.Artist,
	)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error updating album")
	}

	defer db.Close()
	return nil
}

func (a *App) TransferAlbum(album Album) error {
	srcDir := filepath.Join(a.config.CollectionPath, album.Artist, album.Album)
	destDir := filepath.Join(a.config.DevicePath, album.Artist, album.Album)

	os.MkdirAll(destDir, 0750)

	songs, err := os.ReadDir(srcDir)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error reading album directory")
	}

	for _, song := range songs {
		if _, err := os.Stat(filepath.Join(destDir, song.Name())); os.IsNotExist(err) {
			if CheckFileExtension(song.Name()) {
				songPath := filepath.Join(srcDir, song.Name())
				err := cp.Copy(songPath, filepath.Join(destDir, song.Name()))
				if err != nil {
					fmt.Println(err)
					return errors.New("Error copying files")
				}
			}
		}
	}
	return nil
}

func (a *App) RemoveAlbumsFromDevice(ids []int) error {
	// would like to improve the implementation of this...
	// ids should be a slice containing the ids of albums on the device
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error connecting to database")
	}

	placeholders := strings.Trim(strings.Repeat("?,", len(ids)), ",")
	query := fmt.Sprintf("UPDATE albums SET IsOnDevice = 0 WHERE Id NOT IN (%s)", placeholders)

	stmt, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error querying database")
	}

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error querying database")
	}

	query = fmt.Sprintf("SELECT Album, Artist FROM albums WHERE Id NOT IN (%s)", placeholders)
	stmt, err = db.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error querying database")
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Println(err)
		return errors.New("Error querying database")
	}

	for rows.Next() {
		var album string
		var artist string

		err = rows.Scan(&album, &artist)
		if err != nil {
			fmt.Println(err)
		}

		artistPath := filepath.Join(a.config.DevicePath, artist)
		albumPath := filepath.Join(artistPath, album)
		os.RemoveAll(albumPath)

		if dirs, _ := os.ReadDir(artistPath); len(dirs) == 0 {
			os.RemoveAll(artistPath)
		}
	}
	return nil
}
