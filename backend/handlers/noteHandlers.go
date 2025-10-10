package handlers

import (
	"encoding/json"
	"net/http"
	"noteslord/database"
	"noteslord/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	database.DB.Find(&notes)
	json.NewEncoder(w).Encode(notes)
}

func GetNoteByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(note)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if note.Title == "" || note.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}
	database.DB.Create(&note)
	json.NewEncoder(w).Encode(note)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
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
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)
	database.DB.Delete(&models.Note{}, idInt)
	w.WriteHeader(http.StatusNoContent)
}
