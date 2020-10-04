package cmd

import (
	"bufio"
	"fmt"
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
			// viper.Set("API-Key", key)
			// fmt.Println(viper.Get("API-Key"))
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
	cobra.OnInitialize(initConfig)

	// viper.Debug()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.yaml", "config file (default is $HOME/.clockify-cli/config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	//root.go
	rootCmd.AddCommand(addKeyCmd)
	rootCmd.AddCommand(resetViperCmd)
	// project.go
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(offProjectsCmd)
	// entry.go
	rootCmd.AddCommand(startActivityCmd)
	rootCmd.AddCommand(stopActivityCmd)

	startActivityCmd.Flags().BoolVarP(&defFlag, "default", "d", false, "Use default project id.")
	viper.BindPFlag("default", startActivityCmd.Flags().Lookup("default-project"))
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
