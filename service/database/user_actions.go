package database

import (
	"database/sql"
	"strings"
	"fmt"
)

type User struct {
	Username  string
	Following string
	Banned    string
}


func (db *appdbimpl) CheckUsername(username string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Users WHERE username = ?)"
	err := db.c.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (db *appdbimpl) AddUser(username string) error {
	_, err := db.c.Exec("INSERT INTO Users (username) VALUES (?)", username)
	return err
}

func (db *appdbimpl) UpdateUsername(oldUsername, newUsername string) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Step 1: Update the username in the Users table
	_, err = tx.Exec("UPDATE Users SET username = ? WHERE username = ?", newUsername, oldUsername)
	if err != nil {
		return err
	}

	// Step 2: Update the username in the Images table
	_, _ = tx.Exec("UPDATE Images SET username = ? WHERE username = ?", newUsername, oldUsername)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetUser(username string) (User, error) {
	var user User

	// Execute the SELECT query to retrieve user data
	var following sql.NullString
	var banned sql.NullString
	err := db.c.QueryRow("SELECT username, following, banned FROM Users WHERE username = ?", username).
		Scan(&user.Username, &following, &banned)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found, return an empty user and a nil error
			return User{}, fmt.Errorf("user not found")
		}
		// Other error occurred, return the error
		return User{}, err
	}

	// Set the 'Following' field of the user struct
	if following.Valid {
		user.Following = following.String
	} else {
		user.Following = "" // or any default value you want to use
	}

	// Set the 'Banned' field of the user struct
	if banned.Valid {
		user.Banned = banned.String
	} else {
		user.Banned = "" // or any default value you want to use
	}

	return user, nil
}

func (db *appdbimpl) FollowUsername(username, followingusername string) error {
	// Check if the user exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	// Fetch the current following string
	var following string
	err = db.c.QueryRow("SELECT following FROM Users WHERE username = ?", username).Scan(&following)
	if err != nil {
		return err
	}

	// Append the new followingusername to the existing string (comma-separated)
	if following != "" {
		following += ","
	}
	following += followingusername
	// Update the 'following' column for the user
	_, err = db.c.Exec("UPDATE Users SET following = ? WHERE username = ?", following, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnfollowUsername(username, unfollowingusername string) error {
	// Check if the user exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	// Fetch the current following string
	var following string
	err = db.c.QueryRow("SELECT following FROM Users WHERE username = ?", username).Scan(&following)
	if err != nil {
		return err
	}

	// Split the following string into a slice of usernames
	followingList := strings.Split(following, ",")

	// Find and remove the unfollowingusername from the slice
	var updatedFollowingList []string
	for _, user := range followingList {
		if user != unfollowingusername {
			updatedFollowingList = append(updatedFollowingList, user)
		}
	}

	// Join the updated slice back into a comma-separated string
	updatedFollowing := strings.Join(updatedFollowingList, ",")

	// Update the 'following' column for the user
	_, err = db.c.Exec("UPDATE Users SET following = ? WHERE username = ?", updatedFollowing, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) BanUsername(username, banusername string) error {
	// Check if the user exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	var banned string
	err = db.c.QueryRow("SELECT banned FROM Users WHERE username = ?", username).Scan(&banned)
	if err != nil {
		return err
	}

	if banned != "" {
		banned += ","
	}
	banned += banusername
	_, err = db.c.Exec("UPDATE Users SET banned = ? WHERE username = ?", banned, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnbanUsername(username, unbanusername string) error {

	// Check if the user exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	var banned string
	err = db.c.QueryRow("SELECT banned FROM Users WHERE username = ?", username).Scan(&banned)
	if err != nil {
		return err
	}

	bannedList := strings.Split(banned, ",")

	var updatedBannedList []string
	for _, user := range bannedList {
		if user != unbanusername {
			updatedBannedList = append(updatedBannedList, user)
		}
	}

	// Join the updated slice back into a comma-separated string
	updatedBanned := strings.Join(updatedBannedList, ",")

	// Update the 'following' column for the user
	_, err = db.c.Exec("UPDATE Users SET banned = ? WHERE username = ?", updatedBanned, username)
	if err != nil {
		return err
	}

	return nil
}
