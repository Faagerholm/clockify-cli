package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type project struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List workspace projects",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		workspace := viper.Get("workspace")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/projects", workspace), nil)
		req.Header.Set("X-Api-Key", key)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		results := []project{}
		jsonErr := json.Unmarshal(body, &results)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		for i, r := range results {
			if i%3 == 0 {
				fmt.Println("")
			}
			fmt.Printf("(%d) %s, %s.", i+1, r.Name, r.Id)
			for i := len(r.Name) / 17; i < 4; i++ {
				fmt.Print("\t")
			}
		}
		fmt.Print("\n\nSave a project as default: ")
		reader := bufio.NewReader(os.Stdin)
		value, err := reader.ReadString('\n')
		if err == nil {
			l := strings.Trim(value, "\n")
			fmt.Println(l)
			v, _ := strconv.Atoi(l)
			if err != nil {
				log.Fatal(err)
			} else {
				p := results[v-1]
				viper.Set("default-project", p.Id)
				viper.WriteConfig()
				fmt.Println("Default project set:", p)
			}

		}
	},
}
