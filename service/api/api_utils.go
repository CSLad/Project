package api

import (
	"clean/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, _ := rt.db.CheckUsername(requestBody.Username)

	if exists {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"Username": requestBody.Username,
			"Message":  "Successful login into existing account",
		})
		return
	}

	if err := rt.db.AddUser(requestBody.Username); err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"Username": requestBody.Username,
		"Message":  "Successful sign up and login",
	})
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into a struct
	var requestBody struct {
		username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := ps.ByName("username")

	if err := rt.db.UpdateUsername(username, requestBody.username); err != nil {
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := ps.ByName("username")

	if err := rt.db.FollowUsername(username, requestBody.Username); err != nil {
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := ps.ByName("username")
	if err := rt.db.UnfollowUsername(username, requestBody.Username); err != nil {
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := ps.ByName("username")
	if err := rt.db.BanUsername(username, requestBody.Username); err != nil {

		http.Error(w, "Failed to ban user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := ps.ByName("username")
	if err := rt.db.UnbanUsername(username, requestBody.Username); err != nil {
		http.Error(w, "Failed to unban user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	username := ps.ByName("username")
	if username == "" {
		http.Error(w, "Missing username in path", http.StatusBadRequest)
		return
	}

	images, err := rt.db.GetStream(username)
	if err != nil {
		http.Error(w, "Failed to retrieve stream", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(images); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) uploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var requestBody struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	imageURL := ps.ByName("imageurl")
	if err := rt.db.InsertImage(imageURL, requestBody.Username); err != nil {
		http.Error(w, "Failed to insert image into the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	imageURL := ps.ByName("imageurl")
	if imageURL == "" {
		http.Error(w, "Missing image URL in path", http.StatusBadRequest)
		return
	}

	if err := rt.db.RemoveImage(imageURL); err != nil {
		http.Error(w, "Failed to delete image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	imageURL := ps.ByName("imageurl")
	if err := rt.db.AddLike(imageURL); err != nil {
		http.Error(w, "Failed to add like to the image", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) addComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var requestBody struct {
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	imageURL := ps.ByName("imageurl")
	if err := rt.db.AddComment(imageURL, requestBody.Comment); err != nil {
		http.Error(w, "Failed to add comment to the image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) removeComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var requestBody struct {
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	imageURL := ps.ByName("imageurl")
	if err := rt.db.RemoveComment(imageURL, requestBody.Comment); err != nil {
		http.Error(w, "Failed to remove comment from the image", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) getImageInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	imageURL := ps.ByName("imageurl")
	if imageURL == "" {
		http.Error(w, "Missing image URL in path", http.StatusBadRequest)
		return
	}

	image, err := rt.db.GetImage(imageURL)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(image); err != nil {
		http.Error(w, "Failed to encode image data", http.StatusInternalServerError)
		return
	}
}
