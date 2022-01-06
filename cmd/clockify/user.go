package clockify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "get current user",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.GetString("API-KEY")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api.clockify.me/api/v1/user", nil)
		req.Header.Set("X-API-KEY", key)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&result)
		if result != nil {
			fmt.Println("Found user:", result["email"], result["id"])

			viper.Set("USER-ID", result["id"])
			viper.Set("WORKSPACE", result["activeWorkspace"])
			fmt.Println("Updating config with user id.")
			viper.WriteConfig()
		} else {
			log.Println("Could not find user.. try again later.")
		}
	},
}
