/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CheckUsername(username string) (bool, error)
	AddUser(username string) error
	UpdateUsername(oldUsername, newUsername string) error
	GetUser(username string) (User, error)
	FollowUsername(username, followingUsername string) error
	UnfollowUsername(username, unfollowingusername string) error
	BanUsername(username, banusername string) error
	UnbanUsername(username, unbanusername string) error
	GetUserPhotos(username string) ([]Image, error)
	
	GetStream(username string) ([]Image, error)
	InsertImage(imageURL, username string) (int64, error)
	RemoveImage(imageID int64) error
	AddLike(imageID int64) error
	RemoveLike(imageID int64) error
	AddComment(imageID int64, comment string) error
	RemoveComment(imageID int64, commentToRemove string) error
	GetImage(imageID int64) (Image, error)
	
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logger.Infof("Loading Table Users")

	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Infof("No TABLE Users, Initializing")
		sqlStmt := `CREATE TABLE Users (
						username TEXT PRIMARY KEY,
						following TEXT,
						banned TEXT
					);`
		_, err := db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// logger.Infof("Inserting test user")
	// _, err = db.Exec(`INSERT INTO Users (username, following, banned) VALUES (?, ?, ?)`,
	// 	"testuser", "alice,bob", "")
	// if err != nil {
	// 	return nil, fmt.Errorf("error inserting test user: %w", err)
	// }

	logger.Infof("Loading Table Images")

	var tableImagesName string
	errs := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Images';`).Scan(&tableImagesName)
	if errors.Is(errs, sql.ErrNoRows) {
		logger.Infof("No TABLE Images, Initializing")
		sqlStmt := `CREATE TABLE Images (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						imageurl TEXT UNIQUE,
						username TEXT,
						likes INTEGER,
						comments TEXT,
						created_at DATETIME
				);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// logger.Infof("Inserting test images")
	// _, err = db.Exec(`INSERT INTO Images (imageurl, username, likes, comments, created_at)
	// 				VALUES
	// 				(?, ?, ?, ?, datetime('now', '-1 hour')),
	// 				(?, ?, ?, ?, datetime('now', '-2 hour'))`,
	// 	"https://i.guim.co.uk/img/media/c8b1d78883dfbe7cd3f112495941ebc0b25d2265/256_0_3840_2304/master/3840.jpg?width=1200&height=1200&quality=85&auto=format&fit=crop&s=3ac0dfad4c2edc23c843d59f1ec08cc8", "alice", 10, "YES",
	// 	"https://img.freepik.com/foto-gratuito/rappresentazione-3d-del-ragazzo-del-fumetto_23-2150797600.jpg?semt=ais_hybrid&w=740", "bob", 5, "LMAO")
	// if err != nil {
	// 	return nil, fmt.Errorf("error inserting test images: %w", err)
	// }

	// if _, err := db.Exec(`DROP TABLE IF EXISTS Images`); err != nil {
	//     return nil, fmt.Errorf("dropping Images table: %w", err)
	// }
	// if _, err := db.Exec(`DROP TABLE IF EXISTS Users`); err != nil {
	//     return nil, fmt.Errorf("dropping Users table: %w", err)
	// }

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
