package api

import (
	"clean/service/api/reqcontext"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

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
		ImageURL string `json:"imageurl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := rt.db.InsertImage(requestBody.ImageURL, requestBody.Username)
	if err != nil {
		http.Error(w, "Failed to insert image into the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"Username": requestBody.Username,
		"imageId":  id,
	})
}

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	idStr := ps.ByName("imageid")
	if idStr == "" {
		http.Error(w, "Missing image id in path", http.StatusBadRequest)
		return
	}

	imageID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid image id", http.StatusBadRequest)
		return
	}

	if err := rt.db.RemoveImage(imageID); err != nil {
		http.Error(w, "Failed to delete image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	idStr := ps.ByName("imageid")
	imageID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid image id")
		http.Error(w, "Invalid image id", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddLike(imageID); err != nil {
		ctx.Logger.WithError(err).Error("Failed to like the image")
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

	idStr := ps.ByName("imageid")
	imageID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid image id", http.StatusBadRequest)
		return
	}
	if err := rt.db.AddComment(imageID, requestBody.Comment); err != nil {
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

	idStr := ps.ByName("imageid")
	imageID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid image id", http.StatusBadRequest)
		return
	}
	if err := rt.db.RemoveComment(imageID, requestBody.Comment); err != nil {
		http.Error(w, "Failed to remove comment from the image", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (rt *_router) getImageInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	idStr := ps.ByName("imageid")
	if idStr == "" {
		http.Error(w, "Missing image id in path", http.StatusBadRequest)
		return
	}

	imageID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid image id", http.StatusBadRequest)
		return
	}

	image, err := rt.db.GetImage(imageID)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(image); err != nil {
		http.Error(w, "Failed to encode image data", http.StatusInternalServerError)
		return
	}
}
