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

	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type project struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Select default workspace project",
	Long: `Display all workspace projects and 
	select the default project to use when starting a timer`,
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		workspace := viper.Get("WORKSPACE")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/projects", workspace), nil)
		req.Header.Set("X-API-KEY", key)
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

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 12, 8, 12, '\t', 0)

		for i, r := range results {
			if i > 0 && i%3 == 0 {
				fmt.Fprintf(writer, "\n")
			}
			fmt.Fprintf(writer, "(%d) %s\t", i+1, r.Name)
		}
		writer.Flush()
		fmt.Print("\n\nSave a project as default (number): ")
		reader := bufio.NewReader(os.Stdin)
		value, err := reader.ReadString('\n')
		if err == nil {
			l := strings.Trim(value, "\n")
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

var offProjectsCmd = &cobra.Command{
	Use:   "off-projects",
	Short: "Select which projects should be omitted from reports",
	Long: `Display all projects and select which shouldn't be included in 'saldo' report.
	By not selecting these projects the CLI cannot figure out which projects not to exclude when counting your saldo.`,
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		workspace := viper.Get("WORKSPACE")
		client := &http.Client{}
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/projects", workspace), nil)
		req.Header.Set("X-API-KEY", key)
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

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 12, 8, 12, '\t', 0)

		for i, r := range results {
			if i > 0 && i%3 == 0 {
				fmt.Fprintf(writer, "\n")
			}
			fmt.Fprintf(writer, "(%d) %s\t", i+1, r.Name)
		}
		writer.Flush()
		fmt.Print("\n\nSelect projects to exclude from report (number)\nThe projects should be comma-separated\nAny previous projects will be overwritten: ")
		reader := bufio.NewReader(os.Stdin)
		value, err := reader.ReadString('\n')
		if err == nil {
			l := strings.Trim(value, "\n")
			p := strings.Split(l, ",")
			var ps []string
			for _, s := range p {
				v, _ := strconv.Atoi(strings.TrimSpace(s))
				ps = append(ps, results[v-1].Id)
				fmt.Println(results[v-1])
			}
			if err != nil {
				log.Fatal(err)
			} else {
				viper.Set("off-projects", ps)
				viper.WriteConfig()
				fmt.Println("Off projects updated,", ps)
			}
		}
	},
}
