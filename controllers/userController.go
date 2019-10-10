package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jLemmings/GoCV/models"
	"github.com/jLemmings/GoCV/utils"
	"io/ioutil"
	"log"
	"net/http"
)

var GetUsers = func(w http.ResponseWriter, r *http.Request) {
	data := models.GetUsers()
	resp := utils.Message(true, "success")
	resp["data"] = data
	utils.Respond(w, resp)
}

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Print(err)
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}

	resp := user.Create()
	utils.Respond(w, resp)
}

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	data := models.GetUser(request["id"])
	resp := utils.Message(true, "success")
	resp["data"] = data
	utils.Respond(w, resp)
}

var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)

	id := request["id"]
	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(user)

	utils.HandleErr(err)

	updatedUser := models.UpdateUser(id, *user)
	resp := utils.Message(true, "success")
	resp["data"] = updatedUser
	utils.Respond(w, resp)
}

var UpdateUserClaim = func(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	id := request["id"]
	claims := map[string]interface{}{"admin": true}
	err := models.GetAuth().SetCustomUserClaims(context.Background(), id, claims)
	utils.HandleErr(err)

	user, err := models.GetAuth().GetUser(context.Background(), id)
	utils.HandleErr(err)

	resp := utils.Message(true, "success")
	fmt.Print("Added Claim: ", user.CustomClaims, "to UID: ", user.UID)
	utils.Respond(w, resp)
}

var CreateBackup = func(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	id := request["id"]
	fmt.Println(id)

	fmt.Println("IN CREATE BACKUP")
	results, err := models.GetDB().NewRef("users").OrderByKey().GetOrdered(context.Background())
	if err != nil {
		log.Fatal("RESULTS ERROR: ", err)
	}

	if len(results) != 0 {
		var user models.User
		for _, r := range results {
			err := r.Unmarshal(&user)
			if err != nil {
				log.Fatalln("Error unmarshaling result:", err)
			}
			if user.ID == id {
				fmt.Println("USER WAS FOUND:", user.ID)
			}
		}

		fileName := "temp/" + user.ID + "_backup.bak"

		userByte, err := json.Marshal(user)

		if err != nil {
			fmt.Println("Byte ERROR: ", err)
		}
		err = ioutil.WriteFile(fileName, userByte, 0644)
		if err != nil {
			fmt.Println("File ERROR: ", err)
		}
		http.ServeFile(w, r, fileName)
	} else {
		resp := utils.Message(true, "failed")
		resp["data"] = "NO RESULTS FOUND"
		utils.Respond(w, resp)
	}

}
