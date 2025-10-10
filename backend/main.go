package main

import (
	"encoding/json" //converts Go data structs to JSON format and vice versa
	"fmt"   	 //the format package,for formatted I/O operations
	"log"       //used for logging errors and information to the console
	"net/http"  //creates http server and clients
	"strconv"  //converts strings to other data types and vice versa
	"strings"  //provides functions for manipulating strings
	"gorm.io/driver/sqlite"  //SQLite driver for GORM, allows GORM to communicate with SQLite DB
	"gorm.io/gorm"   //GORM library for ORM (Object-Relational Mapping)
)


type Note struct { //A custom data type to represent a note, a struct in Go
	//the backticks-enclosed contents are called struct tags, are metadata read by libraries like GORM and encoding/json
	//the json tag specifies the key name when the struct is serialized to JSON
	//the gorm tag provides instructions to GORM about how to handle this field in the database
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"` //example, here json:"id" tells the JSON encoder/decoder to use "id" as the key for this field
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Global DB connection----declares a global variable db of type *gorm.DB to hold connection to the database
var db *gorm.DB // db is our global GORM database connection

// Initialize SQLite database
func initDB() {
	var err error //initialize an error variable to capture any errors during database connection
	db, err = gorm.Open(sqlite.Open("notes.db"), &gorm.Config{}) //opens a connection to a SQLite database file named notes.db
	if err != nil {
		log.Fatal("Failed to connect to database:", err) //if there's an error the program logs the error and exits
	}

	// Auto-migrate Note model (creates table if not exists)
	if err := db.AutoMigrate(&Note{}); err != nil { //.AutoMigrate() is a GORM method that automatically creates or updates the database schema to match the Note struct
		// &Note{} is a pointer of an instance of the Note struct, telling GORM to create a table based on this struct
		//when db.AutoMigrate(&Note{}) is called, GORM checks if a table for the Note struct exists in the database
		//if it doesn't exist, GORM creates it with columns corresponding to the fields in the Note struct
		//if it does exist, GORM ensures that the table schema matches the current definition of the Note struct, adding any missing columns as needed
		log.Fatal("Failed to migrate database:", err)
	}
     //The AutoMigrate method returns one value: an *gorm.DB instance, which includes an internal Error field. 
	 // In Go, ORM libraries often use this chained pattern. The actual error is retrieved by accessing the .Error property of the returned object.
	fmt.Println("Connected to SQLite database and migrated Note model")
}

// Http Handlers, api endpoints, the below handlwers handle incoming HTTP requests
func getNotesHandler(w http.ResponseWriter, r *http.Request) { //2 arguments: w of type http.ResponseWriter, used to send responses back to the client
	//r of type *http.Request, represents the incoming HTTP request
	//This handler fetches all notes from the database and returns them as a JSON array
	w.Header().Set("Content-Type", "application/json") //this line sets the Content-Type header of the HTTP response to application/json,
	// indicating that the response body will contain JSON data
	var notes []Note //declares a slice of Note structs to hold the fetched notes, slice is like a dynamic array
	//db.Find(&notes) is a GORM method that retrieves all records from the notes table and populates the notes slice with them
	//the &notes is a pointer to the notes slice, allowing GORM to modify it directly
	//if there's an error during the database query, it sends an HTTP 500 error response back to the client
	//if successful, it encodes the notes slice as JSON and writes it to the HTTP response body
	if err := db.Find(&notes).Error; err != nil {
		http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notes)
}

// Handler: Create a new note
func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newNote Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := db.Create(&newNote).Error; err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

// Handler: Update note by ID
func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var note Note
	if err := db.First(&note, id).Error; err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	var updatedData Note
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	note.Title = updatedData.Title
	note.Content = updatedData.Content

	if err := db.Save(&note).Error; err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}

// Handler: Delete note by ID
func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&Note{}, id).Error; err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()

	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getNotesHandler(w, r)
		case http.MethodPost:
			createNoteHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			updateNoteHandler(w, r)
		case http.MethodDelete:
			deleteNoteHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
