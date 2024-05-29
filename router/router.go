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

	// Create a subrouter for protected user routes
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.Use(middlewares.JWTAuth)
	userRouter.HandleFunc("/{userId}", controllers.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/{userId}", controllers.DeleteUser).Methods("DELETE")

	// Public Photo routes
	r.HandleFunc("/photos", controllers.GetPhotos).Methods("GET")

	// Protected Photo routes with JWT
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.JWTAuth)
	api.HandleFunc("/photos", controllers.CreatePhoto).Methods("POST")
	api.HandleFunc("/photos/{photoId}", controllers.UpdatePhoto).Methods("PUT")
	api.HandleFunc("/photos/{photoId}", controllers.DeletePhoto).Methods("DELETE")
	api.HandleFunc("/photos/profile", controllers.SetProfilePhoto).Methods("POST")
	api.HandleFunc("/photos/profile", controllers.GetProfilePhoto).Methods("GET")

	return r
}
