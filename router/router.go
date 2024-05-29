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

	// Protected Photo routes with JWT
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.JWTAuth)
	api.HandleFunc("/photos/profile", controllers.GetProfilePhoto).Methods("GET")  // Mendapatkan foto profil
	api.HandleFunc("/photos/profile", controllers.SetProfilePhoto).Methods("POST") // Mengatur sebuah foto sebagai foto profil
	api.HandleFunc("/photos/{photoId}", controllers.DeletePhoto).Methods("DELETE") // Menghapus sebuah foto

	return r
}
