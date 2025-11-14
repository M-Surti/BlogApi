package main

import (
	"BlogApi/database"
	"BlogApi/handlers"
	"BlogApi/middleware"

	"github.com/gorilla/mux"

	"fmt"

	"net/http"
)


func main() {
	database.InitDB()


	//public routes

	r := mux.NewRouter()


	r.HandleFunc("/register", handlers.RegisterUser)
	r.HandleFunc("/login", handlers.Login)
	r.HandleFunc("/posts", handlers.GetALLPost)

	// protected routes
	protectedroute := r.PathPrefix("/").Subrouter()
	protectedroute.Use(middleware.AuthMiddleware)

	protectedroute.HandleFunc("/create", handlers.CreatePost)

	//start the server
	fmt.Println("server is running")
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}

