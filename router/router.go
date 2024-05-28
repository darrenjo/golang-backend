package router

import (
	"backend-api/controllers"
	"backend-api/middlewares"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/users/register", controllers.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", controllers.LoginUser).Methods("POST")
	r.HandleFunc("/users/{userId}", controllers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{userId}", controllers.DeleteUser).Methods("DELETE")

	// Public Photo routes
	r.HandleFunc("/photos", controllers.GetPhotos).Methods("GET")

	// Protected Photo routes with JWT
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.JWTAuth)
	api.HandleFunc("/photos", controllers.CreatePhoto).Methods("POST")
	api.HandleFunc("/photos/{photoId}", controllers.UpdatePhoto).Methods("PUT")
	api.HandleFunc("/photos/{photoId}", controllers.DeletePhoto).Methods("DELETE")

	return r
}
