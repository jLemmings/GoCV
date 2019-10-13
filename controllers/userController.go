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

		/**

		fileName := "temp/" + user.ID + "_backup.bak"

		userByte, err := json.Marshal(user)

		if err != nil {
			fmt.Println("Byte ERROR: ", err)
		}
		err = ioutil.WriteFile(fileName, userByte, 0644)
		if err != nil {
			fmt.Println("File ERROR: ", err)
		}

		downloadBytes, err := ioutil.ReadFile(fileName)
		fileSize := len(string(downloadBytes))
		mime := http.DetectContentType(downloadBytes)

		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Content-Length", strconv.Itoa(fileSize))
		w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

		http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(downloadBytes))
		*/

		resp := utils.Message(true, "success")
		resp["data"] = user
		utils.Respond(w, resp)
	} else {
		resp := utils.Message(true, "failed")
		resp["data"] = "NO RESULTS FOUND"
		utils.Respond(w, resp)
	}

}

var RestoreBackup = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		panic(err)
	}

	err = models.GetDB().NewRef("users").Set(context.Background(), map[string]*models.User{
		user.ID: {
			ID:              user.ID,
			ProfileImageURL: user.ProfileImageURL,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			Password:        user.Password,
			Bio:             user.Bio,
			GithubProfile:   user.GithubProfile,
			Experience:      user.Experience,
			Education:       user.Education,
		},
	})

	if err != nil {
		log.Fatalln("Error setting value:", err)
	}

	resp := utils.Message(true, "success")
	resp["data"] = "Restored backup successfully"
	utils.Respond(w, resp)

}
