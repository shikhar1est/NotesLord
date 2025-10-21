package handlers

import (
	"encoding/json"
	"net/http"
	"noteslord/database"
	"noteslord/models"
	"noteslord/middleware"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var notes []models.Note
	database.DB.Where("user_id = ?", userID).Find(&notes)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func GetNoteByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var note models.Note

	if err := database.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if note.Title == "" || note.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}
	note.UserID = userID 

	database.DB.Create(&note)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var updated models.Note
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.Title = updated.Title
	note.Content = updated.Content
	database.DB.Save(&note)
	json.NewEncoder(w).Encode(note)
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)

	var note models.Note
	if err := database.DB.First(&note, idInt).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	if note.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
    //cleopatra
	database.DB.Delete(&note)
	w.WriteHeader(http.StatusNoContent)
}
