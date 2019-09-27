package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jLemmings/GoCV/models"
	"github.com/jLemmings/GoCV/utils"
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

var DeleteUser = func(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	models.DeleteUser(request["id"])
	resp := utils.Message(true, "success")
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
