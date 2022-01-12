package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	model "github.com/Faagerholm/clockify-cli/pkg/Model"
	utils "github.com/Faagerholm/clockify-cli/pkg/Utils"
)

func Start(start_time time.Time, project string) {

	api_key, workspace := utils.ExtractAPIKeyAndWorkspace()

	start_time_str := start_time.Format("2006-01-02T15:04:05Z")

	reqBody, err := json.Marshal(map[string]string{
		"start":     start_time_str,
		"projectId": project,
	})

	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/time-entries", workspace), bytes.NewBuffer(reqBody))
	req.Header.Set("X-API-KEY", api_key)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Project started")
	}
	defer resp.Body.Close()
}

func Stop(end_time_str string) model.Entry {
	key, workspace, user := utils.ExtractAPIKeyAndWorkspaceAndUserId()

	reqBody, err := json.Marshal(map[string]string{
		"end": end_time_str,
	})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/user/%s/time-entries", workspace, user), bytes.NewBuffer(reqBody))
	req.Header.Set("X-API-KEY", key)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Project stopped")
	}
	defer resp.Body.Close()

	var entry model.Entry
	json.NewDecoder(resp.Body).Decode(&entry)
	return entry
}

func GetProjects() (
	[]model.Project,
	error,
) {
	key, workspace := utils.ExtractAPIKeyAndWorkspace()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/projects", workspace), nil)
	req.Header.Set("X-API-KEY", key)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	results := []model.Project{}
	jsonErr := json.Unmarshal(body, &results)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return results, nil
}

func GetUser() *model.User {
	key := utils.ExtractAPIKey()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.clockify.me/api/v1/user", nil)
	req.Header.Set("X-API-KEY", key)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var user *model.User

	json.NewDecoder(resp.Body).Decode(&user)
	return user
}

func AddDescription(entryID string, updateEntry model.UpdateEntry) {
	key, workspace := utils.ExtractAPIKeyAndWorkspace()

	reqBody, err := json.Marshal(updateEntry)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/time-entries/%s", workspace, entryID), bytes.NewBuffer(reqBody))
	req.Header.Set("X-API-KEY", key)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	defer resp.Body.Close()
}
