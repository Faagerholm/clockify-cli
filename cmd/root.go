package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	home = os.Getenv("HOME")

	cfgPath string
	// Used for flags.
	cfgFile string
	defFlag bool
	rootCmd = &cobra.Command{
		Use:   "clockify-cli",
		Short: "A Clockify-cli",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute empty command
		},
	}
)

var initCmd = &cobra.Command{
	Use: "init",
	Run: func(cmd *cobra.Command, args []string) {
		key := ""
		reader := bufio.NewReader(os.Stdin)
		if len(args) == 0 {
			fmt.Print("Please enter your api-key:")
			in, _ := reader.ReadString('\n')
			key = strings.Trim(in, "\n")
		} else {
			key = args[0]
		}

		fmt.Print("Are you sure you want to add a new key (Y/N): ")
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		switch char {
		case 'Y':
			viper.Set("API-KEY", key)
			fmt.Println("Saving", viper.Get("API-KEY"), `as your user key, this can be changed later by initializing the same command.
as of now, no more the one key can be used at the same time.`)
			err := viper.WriteConfig() // Find and read the config file
			if err != nil {            // Handle errors reading the config file
				panic(fmt.Errorf("Fatal error config file: %s \n", err))
			}

		case 'N':
			fmt.Println("The key was NOT added.")
		}

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
			fmt.Println("You're ready to go.. check help for more commands")
			viper.WriteConfig()
		} else {
			log.Println("Could not find user.. try again later.")
		}
	},
}

var addKeyCmd = &cobra.Command{
	Use:   "add-key [API-KEY]",
	Short: "Add users API-KEY",
	Long: `Add users API-KEY, get the key from clockify.me/user/settings.
	At the bottom of the page, generate KEY.`,
	Run: func(cmd *cobra.Command, args []string) {

		key := ""
		reader := bufio.NewReader(os.Stdin)
		if len(args) == 0 {
			fmt.Print("Please enter your api-key:")
			in, _ := reader.ReadString('\n')
			key = strings.Trim(in, "\n")
		} else {
			key = args[0]
		}

		fmt.Print("Are you sure you want to add a new key (Y/N): ")
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		switch char {
		case 'Y':
			viper.Set("API-KEY", key)
			fmt.Println("Saving", viper.Get("API-KEY"), `as your user key, this can be changed later by initializing the same command.
as of now, no more the one key can be used at the same time.`)
			err := viper.WriteConfig() // Find and read the config file
			if err != nil {            // Handle errors reading the config file
				panic(fmt.Errorf("Fatal error config file: %s \n", err))
			}

		case 'N':
			fmt.Println("The key was NOT added.")
		}

	},
}
var resetViperCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets viper values",
	Run: func(cmd *cobra.Command, args []string) {
		viper.Reset()
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {

	home, _ := homedir.Dir()
	cobra.OnInitialize(initConfig)

	// viper.Debug()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", home+"/.clockify-cli/config.yaml", "config file (default is $HOME/.clockify-cli/config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	//root.go
	rootCmd.AddCommand(addKeyCmd)
	rootCmd.AddCommand(resetViperCmd)
	rootCmd.AddCommand(initCmd)
	// project.go
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(offProjectsCmd)
	// entry.go
	rootCmd.AddCommand(startActivityCmd)
	rootCmd.AddCommand(stopActivityCmd)

	startActivityCmd.Flags().BoolVarP(&defFlag, "default", "d", false, "Use default project id.")
	// viper.BindPFlag("default", startActivityCmd.Flags().Lookup("default-project"))
	// user.go
	rootCmd.AddCommand(userCmd)
	// workspace.go
	rootCmd.AddCommand(workspaceCmd)
	// utils.go
	rootCmd.AddCommand(versionCmd)
	// report.go
	rootCmd.AddCommand(balanceCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}
		viper.AddConfigPath(home + "/.clockify-cli")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
