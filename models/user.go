package models

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/jLemmings/GoCV/utils"
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

	fmt.Println(user)

	GetDB().Create(user)

	resp := utils.Message(true, "success")
	resp["user"] = user
	return resp
}

func GetUsers() []*User {
	users := make([]*User, 0)
	err := GetDB().Table("users").Find(&users).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return users
}

func GetUser(userId string) *User {
	users := &User{}
	err := GetDB().Table("users").Where("id = ?", userId).First(&users).Error

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return users
}

func UpdateUser(id string, userToUpdate User) *User {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", id).First(&user).Error

	if userToUpdate.FirstName != "" {
		user.FirstName = userToUpdate.FirstName
	}

	if userToUpdate.LastName != "" {
		user.LastName = userToUpdate.LastName
	}
	if userToUpdate.Bio != "" {
		user.Bio = userToUpdate.Bio
	}

	GetDB().Save(&user)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

func DeleteUser(userId string) {
	err := GetDB().Where("id = ?", userId).Delete(&User{})
	if err != nil {
		fmt.Println(err)
	}
}
