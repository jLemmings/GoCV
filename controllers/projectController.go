package controllers

import (
	"context"
	"github.com/google/go-github/v28/github"
	"github.com/gorilla/mux"
	"github.com/jLemmings/GoCV/models"
	"github.com/jLemmings/GoCV/utils"
	"net/http"
)

// https://api.github.com/users/jLemmings/repos
var GetProjects = func(w http.ResponseWriter, r *http.Request) {
	request := mux.Vars(r)
	user := models.GetUser(request["id"])

	opt := &github.RepositoryListOptions{Type: "public"}

	var allRepos []*github.Repository
	for {
		repos, resp, err := models.GetGitClient().Repositories.List(context.Background(), user.GithubProfile, opt)
		utils.HandleErr(err)
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var responseRepos []models.Project
	for _, repo := range allRepos {

		var project models.Project

		project.Name = *repo.Name
		project.URL = *repo.HTMLURL
		project.LastUpdate = *repo.UpdatedAt

		if repo.Description != nil {
			project.Name = *repo.Name
		}

		if repo.Language != nil {
			project.Language = *repo.Language
		}

		responseRepos = append(responseRepos, project)

	}
	resp := utils.Message(true, "success")
	resp["data"] = responseRepos
	utils.Respond(w, resp)
}
