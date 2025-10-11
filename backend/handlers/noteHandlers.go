package handlers

import (
	"encoding/json"
	"net/http"
	"noteslord/database"
	"noteslord/models"
	"strconv"
	"github.com/go-chi/chi/v5"
)


func GetAllNotes(w http.ResponseWriter, r *http.Request) { //this function fetches all notes from the db and returns them as a JSON array
	var notes []models.Note //declares an empty slice of Note structs to hold the fetched notes
	database.DB.Find(&notes) //here we are using GORM's Find method to retrieve all records from the notes table
	// and populate the notes slice with them
	w.Header().Set("Content-Type", "application/json") //this line sets the Content-Type header of the HTTP response to application/json,
	// indicating that the response body will contain JSON data
	json.NewEncoder(w).Encode(notes) //this line takes a Go data  structure (the notes slice) and encodes it as JSON,
	// then writes that JSON data to the HTTP response body
}

func GetNoteByID(w http.ResponseWriter, r *http.Request) { //this function fetches a single note by its ID from the database
	id := chi.URLParam(r, "id") //here we extract the note ID from the URL path using chi.URLParam
	// example URL: /notes/1, id will be "1"
	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil { //tries to find the first record in notes table
		//which is matching the given ID,if it's found then it stores it in the note variable
		//if not found, it returns an error
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(note)
}

func CreateNote(w http.ResponseWriter, r *http.Request) { //this function creates a new note in the database
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil { //reads the JSON data from the request body
		//and decodes it into the note struct, if there's an error during decoding
		//it sends an HTTP 400 error response back to the client
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if note.Title == "" || note.Content == "" { //if any one of the required fields are missing
		//then we return an error
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}
	database.DB.Create(&note) //here we use GORM's Create method to insert the new note into the notes table
	w.WriteHeader(http.StatusCreated) //this line sets the HTTP status code of the response to 201 Created,
	//indicating that a new resource has been successfully created
	json.NewEncoder(w).Encode(note)
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //extract the note ID from the URL path using chi.URLParam
	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil { //we first try to find the existing note by its ID
		//if it's not found, we return a 404 Not Found error
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	var updated models.Note
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil { //we decode the incoming JSON data from the request body into
		// an updated Note struct
		//if there's an error during decoding, we return a 400 Bad Request error
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
	idInt, _ := strconv.Atoi(id) //convert the id string to an integer, also ignoring any conversion error by using '_'
	//gorm's Delete method to remove the note with the specified ID from the database
	//&models.Note{} is a pointer to an empty Note struct, which tells GORM which table to delete from
	//idInt is the ID of the note to delete
	database.DB.Delete(&models.Note{}, idInt)
	w.WriteHeader(http.StatusNoContent)
}
