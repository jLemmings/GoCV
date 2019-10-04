package main

import (
	"flag"
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
	firstName := flag.String("firstName", "", "Your first name [STRING]")
	lastName := flag.String("lastName", "", "Your last name [STRING]")
	email := flag.String("email", "", "Your email [STRING]")
	password := flag.String("password", "", "Your password [STRING]")
	github := flag.String("github", "", "Your GitHub Profile [STRING]")
	flag.Parse()

	if os.Getenv("ENV") != "PROD" {
		err := godotenv.Load()
		*firstName = os.Getenv("firstName")
		*lastName = os.Getenv("lastName")
		*email = os.Getenv("email")
		*password = os.Getenv("password")
		*github = os.Getenv("github")

		if err != nil {
			log.Println("Error loading .env.production file")
		}
	}

	models.GetDB().AutoMigrate(
		&models.User{},
		&models.Experience{},
		&models.Education{},
	)

	users := models.GetDB().First(&models.User{})
	if users.Error != nil {
		if *firstName == "" || *lastName == "" || *email == "" || *password == "" || *github == "" {
			log.Fatal("For setup please enter your user information.")
		}

		models.InitializeDB(*firstName, *lastName, *email, *password, *github)
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
	protectedRoutes.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")
	protectedRoutes.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")

	protectedRoutes.Use(middleware.AuthMiddleware)

	log.Println("Server Started localhost:" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
