package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jLemmings/GoCV/controllers"
	"github.com/jLemmings/GoCV/middleware"
	"github.com/jLemmings/GoCV/models"
	"github.com/jLemmings/GoCV/utils"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	resp := utils.Message(true, "success")
	resp["data"] = "Hello World"
	utils.Respond(w, resp)
}

func main() {
	if os.Getenv("ENV") != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env.production file")
		}
	}

	users, err := models.GetDB().NewRef("users").OrderByKey().GetOrdered(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range users {
		var user models.User
		if err := r.Unmarshal(&user); err != nil {
			log.Fatalln("Error unmarshaling result:", err)
		}
		fmt.Println("USER WAS FOUND:", user.FirstName, user.LastName)
	}

	if users == nil {
		fmt.Println(os.Getenv("firstName"))
		if os.Getenv("firstName") == "" || os.Getenv("lastName") == "" || os.Getenv("email") == "" || os.Getenv("password") == "" || os.Getenv("github") == "" {
			log.Fatal("For setup please enter your user information.")
		}

		models.InitializeFirstUser(os.Getenv("firstName"), os.Getenv("lastName"), os.Getenv("email"), os.Getenv("password"), os.Getenv("github"))
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.CORSMethodMiddleware(router))

	router.Use(middleware.LogMiddleware)

	router.HandleFunc("/users/{id}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/projects/{id}", controllers.GetProjects).Methods("GET")
	router.HandleFunc("/", helloWorld).Methods("GET")

	protectedRoutes := router.PathPrefix("").Subrouter()
	protectedRoutes.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	protectedRoutes.HandleFunc("/users/{id}", controllers.UpdateUserClaim).Methods("POST")
	protectedRoutes.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")

	protectedRoutes.Use(middleware.AuthMiddleware)

	log.Println("Server Started localhost:" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
