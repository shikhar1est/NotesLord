//this is the routes package, where we set the routing configuration for note-related endpoints
//basically tells the application which function to call when a specific HTTP request is made to a certain URL path
package routes

import (
	"github.com/go-chi/chi/v5"
	"noteslord/handlers"
)

func RegisterNoteRoutes(r chi.Router) { //here 'r' is a chi.Router instance, which is used to define routes
	//we define a group of routes under the /notes path
	//each route corresponds to a specific HTTP method (GET, POST, PUT, DELETE)
	//and is associated with a handler function from the handlers package
	//these handler functions contain the logic to process the requests and generate responses
	r.Route("/notes", func(r chi.Router) {
		r.Get("/", handlers.GetAllNotes)
		r.Get("/{id}", handlers.GetNoteByID)
		r.Post("/", handlers.CreateNote)
		r.Put("/{id}", handlers.UpdateNote)
		r.Delete("/{id}", handlers.DeleteNote)
	})
}
