package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"noteslord/database"
	"noteslord/routes"
	"noteslord/config"
)

// Initialize SQLite database
func main() {
	config.LoadEnv()
	database.ConnectDB()

	r := chi.NewRouter() //here we create a router instance using chi.NewRouter()
	r.Use(middleware.Logger) //inbuilt middleware that logs the start and end of each request, 
	// along with some useful data about what was requested and how long it took to process.
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to NotesLord"))
	})

	routes.RegisterAuthRoutes(r)
	routes.RegisterNoteRoutes(r)
	if err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
    log.Printf("%s %s\n", method, route)
    return nil
}); err != nil {
    log.Fatal(err)
}

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
