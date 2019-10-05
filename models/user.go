package models

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/jLemmings/GoCV/utils"
	"log"
)

type User struct {
	ID            string `gorm:"primary_key"`
	FirstName     string
	LastName      string
	Email         string
	Password      string `gorm:"-"`
	Bio           string
	GithubProfile string
	Experience    Experience `gorm:"foreignkey:UserID"`
	Education     Education  `gorm:"foreignkey:UserID"`
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
		return &user
	}
	return &User{}
}

func UpdateUser(userId string, userToUpdate User) *User {
	user := &User{}
	err := GetDB().NewRef("users"+userId).Update(context.Background(), map[string]interface{}{
		"FirstName":     userToUpdate.FirstName,
		"LastName":      userToUpdate.LastName,
		"Email":         userToUpdate.Email,
		"Password":      userToUpdate.Password,
		"Bio":           userToUpdate.Password,
		"GithubProfile": userToUpdate.GithubProfile,
		"Experience":    userToUpdate.Experience,
		"Education":     userToUpdate.Education,
	})

	if err != nil {
		log.Println("Error updating user: ", err)
	}
	return user
}

func InitializeFirstUser(firstName string, lastName string, email string, password string, github string) {
	user := User{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Password:      password,
		Bio:           "",
		GithubProfile: github,
		Experience:    Experience{},
		Education:     Education{},
	}

	log.Println(user)
	user.Create()
}
