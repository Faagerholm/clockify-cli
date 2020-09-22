package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Get workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api.clockify.me/api/v1/workspaces", nil)
		req.Header.Set("X-API-KEY", key)
		req.Header.Set("Host", "api.clockify.me")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&result)

		fmt.Println("Found user:", result)
	},
}
