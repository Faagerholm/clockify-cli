package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "get current user",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api.clockify.me/api/v1/user", nil)
		req.Header.Set("X-Api-Key", key)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var result map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&result)

		fmt.Println("Found user:", result["email"], result["id"])

		viper.Set("user_id", result["id"])
		viper.Set("workspace", result["activeWorkspace"])
		viper.WriteConfig()
	},
}

var setUserCmd = &cobra.Command{
	Use:   "add-key [API-KEY]",
	Short: "Add users API-KEY",
	Long: `Add users API-Key, get the key from clockify.me/user/settings.
	At the bottom of the page, generate KEY.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Print("Are you sure you want to add a new key (Y/N): ")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		switch char {
		case 'Y':
			fmt.Println("Saving", args[0], `as your user key, this can be changed later by initializing the same command.
			as of now, no more the one key can be used at the same time.`)

			viper.Set("API-KEY", args[0])
			err := viper.WriteConfig() // Find and read the config file
			if err != nil {            // Handle errors reading the config file
				panic(fmt.Errorf("Fatal error config file: %s \n", err))
			}
		case 'N':
			fmt.Println("The key was not added.")
		}

	},
}

func getConfigs() {
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println(err)
	}
}
