package utils

import (
	"log"

	"github.com/spf13/viper"
)

func ExtractAPIKey() string {
	api_key := viper.GetString("api-key")

	if api_key == "" {
		log.Fatal("No API key found. Please set the API key with 'clockify-cli config set api-key <api-key>'")
	}

	return api_key
}

func ExtractWorksapce() string {
	workspace := viper.GetString("WORKSPACE")

	if workspace == "" {
		log.Fatal("No workspace found. Please set the workspace with 'clockify-cli config set workspace <workspace-id>'")
	}

	return workspace
}

func ExtractAPIKeyAndWorkspace() (string, string) {
	api_key := viper.GetString("api-key")
	workspace := viper.GetString("WORKSPACE")

	if api_key == "" {
		log.Fatal("No API key found. Please set the API key with 'clockify-cli config set api-key <api-key>'")
	} else if workspace == "" {
		log.Fatal("No workspace found. Please set the workspace with 'clockify-cli config set workspace <workspace-id>'")
	}

	return api_key, workspace
}

func ExtractAPIKeyAndWorkspaceAndUserId() (
	api_key string,
	workspace string,
	user_id string,
) {
	api_key, workspace = ExtractAPIKeyAndWorkspace()
	user_id = viper.GetString("USER-ID")

	if user_id == "" {
		log.Fatal("No user id found. Please set the user id with 'clockify-cli config set user-id <user-id>'")
	}

	return api_key, workspace, user_id
}

func IsUserAuthenticated() bool {
	user_id := viper.GetString("USER-ID")

	if user_id == "" {
		return false
	}

	return true
}
