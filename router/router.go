package router

import (
	"backend-api/controllers"
	"backend-api/middlewares"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// User routes (No JWT required)
	r.HandleFunc("/users/register", controllers.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", controllers.LoginUser).Methods("POST")

	// Subrouter for protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.JWTAuth)

	// User routes with JWT protection
	api.HandleFunc("/users/{userId}", controllers.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{userId}", controllers.DeleteUser).Methods("DELETE")

	// Photo routes with JWT protection
	api.HandleFunc("/photos", controllers.CreatePhoto).Methods("POST")
	api.HandleFunc("/photos/{photoId}", controllers.UpdatePhoto).Methods("PUT")
	api.HandleFunc("/photos/{photoId}", controllers.DeletePhoto).Methods("DELETE")
	api.HandleFunc("/photos/profile", controllers.SetProfilePhoto).Methods("POST")
	api.HandleFunc("/photos/profile", controllers.GetProfilePhoto).Methods("GET")

	// Public Photo routes
	r.HandleFunc("/photos", controllers.GetPhotos).Methods("GET")

	return r
}
