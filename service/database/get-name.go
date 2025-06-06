package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName() (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM example_table WHERE id=1").Scan(&name)
	return name, err
}

type User struct {
	Username  string
	Following string
	Banned    string
}

type Image struct {
	ImageURL  string    `json:"imageurl"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

func (db *appdbimpl) CheckUsername(username string) (bool, error) {
	fmt.Print("Checking username")
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM Users WHERE username = ?)"
	err := db.c.QueryRow(query, username).Scan(&exists)
	fmt.Print(exists)
	return exists, err
}

// AddUser adds a new user to the database.
func (db *appdbimpl) AddUser(username string) error {
	fmt.Println("Inserting into users")

	// Initialize empty arrays for following and banned
	_, err := db.c.Exec("INSERT INTO Users (username) VALUES (?)", username)
	fmt.Println(err)
	return err
}

func (db *appdbimpl) UpdateUsername(oldUsername, newUsername string) error {
	fmt.Println("we starting to update")
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Step 1: Update the username in the Users table
	_, err = tx.Exec("UPDATE Users SET username = ? WHERE username = ?", newUsername, oldUsername)
	fmt.Println(err)
	if err != nil {
		return err
	}

	// Step 2: Update the username in the Images table
	_, _ = tx.Exec("UPDATE Images SET username = ? WHERE username = ?", newUsername, oldUsername)
	fmt.Println("We on set in images")
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	fmt.Println("Username updated successfully in both tables.")
	return nil
}

func (db *appdbimpl) GetUser(username string) (User, error) {
	var user User

	// Execute the SELECT query to retrieve user data
	var following sql.NullString
	var banned sql.NullString
	err := db.c.QueryRow("SELECT username, following, banned FROM Users WHERE username = $1", username).
		Scan(&user.Username, &following, &banned)
	fmt.Print(err)
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
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	// Fetch the current following string
	var following string
	err = db.c.QueryRow("SELECT following FROM Users WHERE username = $1", username).Scan(&following)
	if err != nil {
		return err
	}

	// Append the new followingusername to the existing string (comma-separated)
	if following != "" {
		following += ","
	}
	following += followingusername
	fmt.Println(following)
	// Update the 'following' column for the user
	_, err = db.c.Exec("UPDATE Users SET following = $1 WHERE username = $2", following, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnfollowUsername(username, unfollowingusername string) error {
	// Check if the user exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	// Fetch the current following string
	var following string
	err = db.c.QueryRow("SELECT following FROM Users WHERE username = $1", username).Scan(&following)
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
	_, err = db.c.Exec("UPDATE Users SET following = $1 WHERE username = $2", updatedFollowing, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) BanUsername(username, banusername string) error {
	// Check if the user exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	var banned string
	err = db.c.QueryRow("SELECT banned FROM Users WHERE username = $1", username).Scan(&banned)
	if err != nil {
		return err
	}

	if banned != "" {
		banned += ","
	}
	banned += banusername
	_, err = db.c.Exec("UPDATE Users SET banned = $1 WHERE username = $2", banned, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnbanUsername(username, unbanusername string) error {

	// Check if the user exists
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM Users WHERE username = $1", username).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", username)
	}

	var banned string
	err = db.c.QueryRow("SELECT banned FROM Users WHERE username = $1", username).Scan(&banned)
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
	_, err = db.c.Exec("UPDATE Users SET banned = $1 WHERE username = $2", updatedBanned, username)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetStream(username string) ([]Image, error) {
	// Step 1: Retrieve the list of usernames that the given username follows
	var following string
	err := db.c.QueryRow("SELECT following FROM Users WHERE username = $1", username).Scan(&following)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	// Step 2: Query the Images table for images posted by the users that the given username follows
	rows, err := db.c.Query("SELECT imageurl, username, likes, comments, created_at FROM Images WHERE username IN ($1) LIMIT 10", following)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Step 3: Iterate over the query results and build the list of images
	var images []Image
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.ImageURL, &image.Username, &image.Likes, &image.Comments, &image.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func (db *appdbimpl) InsertImage(imageURL, username string) error {
	// Get the current time
	currentTime := time.Now()

	// Execute the INSERT query to insert the image URL into the Images table
	_, err := db.c.Exec("INSERT INTO Images (imageurl, username, likes, comments, created_at) VALUES ($1, $2, $3, $4, $5)", imageURL, username, 0, "", currentTime)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveImage(imageURL string) error {
	// Execute the DELETE query to remove the entry associated with the given image URL
	_, err := db.c.Exec("DELETE FROM Images WHERE imageurl = $1", imageURL)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) AddLike(imageURL string) error {
	// Execute the UPDATE query to increment the number of likes for the corresponding image URL
	_, err := db.c.Exec("UPDATE Images SET likes = likes + 1 WHERE imageurl = $1", imageURL)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveLike(imageURL string) error {
	// Execute the UPDATE query to decrement the number of likes for the corresponding image URL
	_, err := db.c.Exec("UPDATE Images SET likes = likes - 1 WHERE imageurl = $1", imageURL)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) AddComment(imageURL, comment string) error {
	// Retrieve the current comments for the image
	var currentComments string
	err := db.c.QueryRow("SELECT comments FROM Images WHERE imageurl = $1", imageURL).Scan(&currentComments)
	if err != nil {
		return err
	}

	// Concatenate the new comment with the current comments, separated by a special character
	newComments := currentComments + "~" + comment

	// Update the comments for the image
	_, err = db.c.Exec("UPDATE Images SET comments = $1 WHERE imageurl = $2", newComments, imageURL)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) RemoveComment(imageURL, commentToRemove string) error {
	// Retrieve the current comments for the image
	var currentComments string
	err := db.c.QueryRow("SELECT comments FROM Images WHERE imageurl = $1", imageURL).Scan(&currentComments)
	if err != nil {
		return err
	}

	// Split the current comments string into individual comments
	comments := strings.Split(currentComments, "~")

	// Find and remove the comment to be removed
	var updatedComments []string
	for _, c := range comments {
		if c != commentToRemove {
			updatedComments = append(updatedComments, c)
		}
	}

	// Join the remaining comments into a single string
	newComments := strings.Join(updatedComments, "~")

	// Update the comments for the image
	_, err = db.c.Exec("UPDATE Images SET comments = $1 WHERE imageurl = $2", newComments, imageURL)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetImage(imageURL string) (Image, error) {
	// Query the Images table for the image with the given imageURL
	var image Image
	err := db.c.QueryRow("SELECT imageurl, username, likes, comments, created_at FROM Images WHERE imageurl = $1", imageURL).Scan(&image.ImageURL, &image.Username, &image.Likes, &image.Comments, &image.CreatedAt)
	if err != nil {
		return Image{}, err
	}
	return image, nil
}
