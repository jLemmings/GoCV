package models

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/google/go-github/v28/github"
	"github.com/jLemmings/GoCV/utils"
	"log"
	"strings"
	"time"
)

type User struct {
	ID              string
	ProfileImageURL string
	FirstName       string
	LastName        string
	Email           string
	Password        string
	Bio             string
	GithubProfile   string
	Experience      []Experience
	Education       []Education
	Projects        []Project
}

func (user *User) Validate() (map[string]interface{}, bool) {
	if user.Password == "" {
		return utils.Message(false, "Password should be on the payload"), false
	}

	if user.Email == "" {
		return utils.Message(false, "Email should be on the payload"), false
	}

	if user.FirstName == "" {
		return utils.Message(false, "FirstName should be on the payload"), false
	}

	if user.LastName == "" {
		return utils.Message(false, "LastName should be on the payload"), false
	}

	return utils.Message(true, "success"), true
}

func (user *User) Create() map[string]interface{} {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	params := (&auth.UserToCreate{}).
		Email(user.Email).
		Password(user.Password)

	fireUser, err := GetAuth().CreateUser(context.Background(), params)
	utils.HandleErr(err)

	user.ID = fireUser.UID

	user.Experience = []Experience{{
		Title:       "My Title",
		Description: "My Description",
		From:        time.Time{},
		To:          time.Time{},
		Tasks:       []string{"Test", "bob", "Fred"},
	}}

	user.Education = []Education{{
		Title:     "My Title",
		Institute: "My Institute",
		From:      time.Time{},
		To:        time.Time{},
	}}

	err = GetDB().NewRef("users/"+user.ID).Set(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	resp := utils.Message(true, "success")
	resp["user"] = user
	return resp
}

func GetUsers() []*User {
	users := make([]*User, 0)
	err := GetDB().NewRef("users").OrderByValue().Get(context.Background(), &User{})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return users
}

func GetUser(userId string) *User {
	results, err := GetDB().NewRef("users").OrderByKey().EqualTo(userId).GetOrdered(context.Background())
	if err != nil {
		log.Fatalln("Error querying database:", err)
	}

	for _, r := range results {
		var user User
		if err := r.Unmarshal(&user); err != nil {
			log.Fatalln("Error unmarshaling result:", err)
		}
		user.Projects = getProjects(user.GithubProfile)
		return &user
	}

	return &User{}
}

func UpdateUser(userId string, userToUpdate User) *User {
	user := &User{}
	err := GetDB().NewRef("users"+userId).Update(context.Background(), map[string]interface{}{
		"ProfileImageURL": userToUpdate.ProfileImageURL,
		"FirstName":       userToUpdate.FirstName,
		"LastName":        userToUpdate.LastName,
		"Email":           userToUpdate.Email,
		"Password":        userToUpdate.Password,
		"Bio":             userToUpdate.Password,
		"GithubProfile":   userToUpdate.GithubProfile,
		"Experience":      userToUpdate.Experience,
		"Education":       userToUpdate.Education,
	})

	if err != nil {
		log.Println("Error updating user: ", err)
	}
	return user
}

func InitializeFirstUser(firstName string, lastName string, email string, password string, github string) {
	user := User{
		ProfileImageURL: "",
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		Password:        password,
		Bio:             "",
		GithubProfile:   github,
		Experience:      []Experience{},
		Education:       []Education{},
	}

	log.Println(user)
	user.Create()
}

func getProjects(githubProfile string) []Project {
	opt := &github.RepositoryListOptions{Type: "public"}

	var allRepos []*github.Repository

	for {
		repos, resp, err := GetGitClient().Repositories.List(context.Background(), githubProfile, opt)
		utils.HandleErr(err)
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var responseRepos []Project

	for _, repo := range allRepos {

		var project Project
		project.Name = *repo.Name
		project.URL = *repo.HTMLURL
		project.LastUpdate = *repo.UpdatedAt

		if repo.Description != nil {
			project.Stack = strings.Split(*repo.Description, ",")
		}

		if repo.Language != nil {
			project.Language = *repo.Language
		}
		responseRepos = append(responseRepos, project)
	}

	return responseRepos
}
