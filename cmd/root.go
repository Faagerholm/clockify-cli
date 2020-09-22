package cmd

import (
	"bufio"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	home = os.Getenv("HOME")

	// Used for flags.
	cfgFile string
	defFlag bool
	rootCmd = &cobra.Command{
		Use:   "clockify-cli",
		Short: "A Clockify-cli",
	}
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup clockify-cli",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet("API-KEY") {
			fmt.Print("Old API-Key found, use that (Y/N): ")
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println(err)
			}
			if char == 'Y' {
				fmt.Println("Using old API-KEY, you can setup a new key with 'clockify-cli add-key' ")
			} else {
				fmt.Print("Enter new key: ")

				reader = bufio.NewReader(os.Stdin)
				key, err := reader.ReadString('\n')
				if err == nil {
					viper.Set("API-KEY", key)
					fmt.Println("New API-KEY set! Happy clocking..")
					viper.WriteConfig()
				}
			}
		}
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "configFile", "c", home+"/.clockify-cli/config.json", "configuration file, default is "+home+"/.clockify-cli/config.json")

	viper.SetDefault("author", "Jimmy Fagerholm fagerholm.jimmy@gmail.com")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	//root.go
	rootCmd.AddCommand(initCmd)
	// project.go
	rootCmd.AddCommand(projectsCmd)
	// entry.go
	rootCmd.AddCommand(startActivityCmd)
	rootCmd.AddCommand(stopActivityCmd)
	startActivityCmd.Flags().BoolVarP(&defFlag, "default", "d", false, "Use default project id.")
	viper.BindPFlag("default", startActivityCmd.Flags().Lookup("default-project"))
	// user.go
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(setUserCmd)
	// workspace.go
	rootCmd.AddCommand(workspaceCmd)
	// utils.go
	rootCmd.AddCommand(versionCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		fmt.Println()
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home + "/.clockify-cli/")
		viper.SetConfigName(".config.json")
		fmt.Println("config file set!")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// readConfig reads config in .clockify-cli/config per default
func readConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.ReadInConfig() // Find and read the config file
	}
}
