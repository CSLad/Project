package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract username from the path parameters
	username := ps.ByName("username")

	// Retrieve user profile from the database
	user, err := rt.db.GetUser(username)
	if err != nil {
		// Handle database error
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	// Marshal user profile to JSON
	userProfileJSON, err := json.Marshal(user)
	if err != nil {
		// Handle JSON marshaling error
		http.Error(w, "Failed to marshal user profile to JSON", http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userProfileJSON)
	if err != nil {
		// Handle write error
		fmt.Println("Failed to write response:", err)
	}
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract username from the path parameters
	username := ps.ByName("username")

	// Retrieve user's stream from the database
	stream, err := rt.db.GetStream(username)
	fmt.Println(err)
	if err != nil {
		// Handle database error
		http.Error(w, "Failed to retrieve user's stream", http.StatusInternalServerError)
		return
	}

	// Marshal user's stream to JSON
	streamJSON, err := json.Marshal(stream)
	if err != nil {
		// Handle JSON marshaling error
		http.Error(w, "Failed to marshal user's stream to JSON", http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(streamJSON)
	if err != nil {
		// Handle write error
		fmt.Println("Failed to write response:", err)
	}
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the image URL from the path parameters
	imageURL := ps.ByName("imageurl")

	// Remove the image from the database
	if err := rt.db.RemoveImage(imageURL); err != nil {
		http.Error(w, "Failed to remove image from the database", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image removed successfully")
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the image URL from the path parameters
	imageURL := ps.ByName("imageurl")

	// Remove a like from the image in the database
	if err := rt.db.RemoveLike(imageURL); err != nil {
		http.Error(w, "Failed to remove like from the image", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Like removed successfully")
}

func (rt *_router) getImageInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract imageURL from the path parameters
	imageURL := ps.ByName("imageurl")

	// Retrieve image info from the database
	image, err := rt.db.GetImage(imageURL)
	if err != nil {
		// Handle database error
		http.Error(w, "Failed to retrieve image info", http.StatusInternalServerError)
		return
	}

	// Marshal image info to JSON
	imageInfoJSON, err := json.Marshal(image)
	if err != nil {
		// Handle JSON marshaling error
		http.Error(w, "Failed to marshal image info to JSON", http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(imageInfoJSON)
	if err != nil {
		// Handle write error
		fmt.Println("Failed to write response:", err)
	}
}
