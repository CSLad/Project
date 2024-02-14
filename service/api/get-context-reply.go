package api

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"fmt"
)

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) getContextReply(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")

    // Check if the request method is OPTIONS (preflight request)
    if r.Method == http.MethodOptions {
        // Respond with allowed methods and headers for CORS preflight
        w.Header().Set("Access-Control-Allow-Methods", "POST")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Check if the request method is POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Decode the request body into a struct
    var requestBody struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Check if the username exists
    exists, _ := rt.db.CheckUsername(requestBody.Username)

    if exists {
        // Username already exists, return 200 OK
		_, _ = w.Write([]byte("Successful login into existing account"))
        w.WriteHeader(http.StatusOK)
        return
    }

    // Username does not exist, add it
    if err := rt.db.AddUser(requestBody.Username); err != nil {
        // Error adding user, return internal server error
        http.Error(w, "Failed to add user", http.StatusInternalServerError)
        return
    }

    // User added successfully, return 201 Created
	_, _ = w.Write([]byte("Successful sign up and login"))
    w.WriteHeader(http.StatusCreated)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")

    // Check if the request method is OPTIONS (preflight request)
    if r.Method == http.MethodOptions {
        // Respond with allowed methods and headers for CORS preflight
        w.Header().Set("Access-Control-Allow-Methods", "PUT")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Check if the request method is PUT
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the username from the path parameter
    username := ps.ByName("username")

    // Decode the request body into a struct
    var requestBody struct {
        username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Update the username
    if err := rt.db.UpdateUsername(username, requestBody.username); err != nil {
        // Error updating username, return internal server error
		fmt.Println(err)
        http.Error(w, "Failed to update username", http.StatusInternalServerError)
        return
    }

    // Username updated successfully, return 200 OK
    w.WriteHeader(http.StatusOK)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")

    // Check if the request method is OPTIONS (preflight request)
    if r.Method == http.MethodOptions {
        // Respond with allowed methods and headers for CORS preflight
        w.Header().Set("Access-Control-Allow-Methods", "PUT")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Check if the request method is PUT
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the username from the path parameter
    username := ps.ByName("username")

    var requestBody struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
	fmt.Println(username,"and", requestBody.Username)
    if err := rt.db.FollowUsername(username, requestBody.Username); err != nil {
		fmt.Println(err)
        http.Error(w, "Failed to follow user", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")

    // Check if the request method is OPTIONS (preflight request)
    if r.Method == http.MethodOptions {
        // Respond with allowed methods and headers for CORS preflight
        w.Header().Set("Access-Control-Allow-Methods", "DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Check if the request method is DELETE
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the username from the path parameter
    username := ps.ByName("username")

    var requestBody struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
	fmt.Println(username,"and", requestBody.Username)
    if err := rt.db.UnfollowUsername(username, requestBody.Username); err != nil {
		fmt.Println(err)
        http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")

    // Check if the request method is OPTIONS (preflight request)
    if r.Method == http.MethodOptions {
        // Respond with allowed methods and headers for CORS preflight
        w.Header().Set("Access-Control-Allow-Methods", "PUT")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Check if the request method is PUT
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the username from the path parameter
    username := ps.ByName("username")

    var requestBody struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
	fmt.Println(username,"and", requestBody.Username)
    if err := rt.db.BanUsername(username, requestBody.Username); err != nil {
		fmt.Println(err)
        http.Error(w, "Failed to ban user", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    w.Header().Set("Content-Type", "application/json")

    // Check if the request method is OPTIONS (preflight request)
    if r.Method == http.MethodOptions {
        // Respond with allowed methods and headers for CORS preflight
        w.Header().Set("Access-Control-Allow-Methods", "DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }

    // Check if the request method is DELETE
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the username from the path parameter
    username := ps.ByName("username")

    var requestBody struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
	fmt.Println(username,"and", requestBody.Username)
    if err := rt.db.UnbanUsername(username, requestBody.Username); err != nil {
		fmt.Println(err)
        http.Error(w, "Failed to unban user", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}

func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Extract the image URL from the path parameters
    imageURL := ps.ByName("imageurl")

    // Decode the request body to get the username
    var requestBody struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Insert the image into the database
    if err := rt.db.InsertImage(imageURL, requestBody.Username); err != nil {
        http.Error(w, "Failed to insert image into the database", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Image inserted successfully")
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Extract the image URL from the path parameters
    imageURL := ps.ByName("imageurl")

    // Add a like to the image in the database
    if err := rt.db.AddLike(imageURL); err != nil {
        http.Error(w, "Failed to add like to the image", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Like added successfully")
}

func (rt *_router) addComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Extract the image URL from the path parameters
    imageURL := ps.ByName("imageurl")

    // Decode the request body to get the comment
    var requestBody struct {
        Comment string `json:"comment"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Add the comment to the image in the database
    if err := rt.db.AddComment(imageURL, requestBody.Comment); err != nil {
        http.Error(w, "Failed to add comment to the image", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Comment added successfully")
}

func (rt *_router) removeComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // Extract the image URL from the path parameters
    imageURL := ps.ByName("imageurl")

    // Decode the request body to get the comment to remove
    var requestBody struct {
        Comment string `json:"comment"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Remove the comment from the image in the database
    if err := rt.db.RemoveComment(imageURL, requestBody.Comment); err != nil {
        http.Error(w, "Failed to remove comment from the image", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Comment removed successfully")
}
