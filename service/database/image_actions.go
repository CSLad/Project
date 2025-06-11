package database

import (
	"fmt"
	"strings"
	"time"
)

type Image struct {
	ID        int64     `json:"id"`
	ImageURL  string    `json:"imageurl"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

func (db *appdbimpl) GetStream(username string) ([]Image, error) {
	var following string
	err := db.c.QueryRow("SELECT following FROM Users WHERE username = ?", username).Scan(&following)
	if err != nil {
		return nil, err
	}

	// Split and clean the following list
	usernames := strings.Split(following, ",")
	var cleaned []string
	for _, u := range usernames {
		u = strings.TrimSpace(u)
		if u != "" {
			cleaned = append(cleaned, u)
		}
	}

	if len(cleaned) == 0 {
		return []Image{}, nil // User follows no one
	}

	// Build placeholders (?, ?, ...) and args
	placeholders := make([]string, len(cleaned))
	args := make([]interface{}, len(cleaned))
	for i, u := range cleaned {
		placeholders[i] = "?"
		args[i] = u
	}

	query := fmt.Sprintf(
		"SELECT id, imageurl, username, likes, comments, created_at FROM Images WHERE username IN (%s) ORDER BY created_at DESC LIMIT 10",
		strings.Join(placeholders, ","),
	)

	rows, err := db.c.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var image Image
		err := rows.Scan(&image.ID, &image.ImageURL, &image.Username, &image.Likes, &image.Comments, &image.CreatedAt)
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

func (db *appdbimpl) InsertImage(imageURL, username string) (int64, error) {
	// Get the current time
	currentTime := time.Now()

	// Execute the INSERT query to insert the image URL into the Images table
	res, err := db.c.Exec("INSERT INTO Images (imageurl, username, likes, comments, created_at) VALUES (?, ?, ?, ?, ?)", imageURL, username, 0, "", currentTime)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *appdbimpl) RemoveImage(imageID int64) error {
	// Execute the DELETE query to remove the entry associated with the given image URL
	_, err := db.c.Exec("DELETE FROM Images WHERE id = ?", imageID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) AddLike(imageID int64) error {
	_, err := db.c.Exec("UPDATE Images SET likes = likes + 1 WHERE id = ?", imageID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveLike(imageID int64) error {
	// Execute the UPDATE query to decrement the number of likes for the corresponding image
	_, err := db.c.Exec("UPDATE Images SET likes = likes - 1 WHERE id = ?", imageID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) AddComment(imageID int64, comment string) error {
	// Retrieve the current comments for the image
	var currentComments string
	err := db.c.QueryRow("SELECT comments FROM Images WHERE id = ?", imageID).Scan(&currentComments)
	if err != nil {
		return err
	}

	// Concatenate the new comment with the current comments, separated by a special character
	newComments := currentComments + "~" + comment

	// Update the comments for the image
	_, err = db.c.Exec("UPDATE Images SET comments = ? WHERE id = ?", newComments, imageID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) RemoveComment(imageID int64, commentToRemove string) error {
	// Retrieve the current comments for the image
	var currentComments string
	err := db.c.QueryRow("SELECT comments FROM Images WHERE id = ?", imageID).Scan(&currentComments)
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
	_, err = db.c.Exec("UPDATE Images SET comments = ? WHERE id = ?", newComments, imageID)
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) GetImage(imageID int64) (Image, error) {
	// Query the Images table for the image with the given ID
	var image Image
	err := db.c.QueryRow("SELECT id, imageurl, username, likes, comments, created_at FROM Images WHERE id = ?", imageID).Scan(&image.ID, &image.ImageURL, &image.Username, &image.Likes, &image.Comments, &image.CreatedAt)
	if err != nil {
		return Image{}, err
	}
	return image, nil
}

