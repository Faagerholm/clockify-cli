package main

import (
	"fmt"
	"os"

	commands "github.com/Faagerholm/clockify-cli/cmd"
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
			// Start menu if no subcommand is given.
			commands.CheckConfigAndPromptSetup()
			commands.Menu()
		},
	}
)

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", home+"/.clockify-cli/config.yaml", "config file (default is $HOME/.clockify-cli/config.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	//root.go
	rootCmd.AddCommand(commands.AddKeyCmd)
	rootCmd.AddCommand(resetViperCmd)
	rootCmd.AddCommand(commands.SetupCmd)
	// menu.go
	rootCmd.AddCommand(commands.MenuCmd)
	// project.go
	rootCmd.AddCommand(commands.DefaultProjectCmd)
	rootCmd.AddCommand(commands.ListProjectsCmd)
	// entry.go
	rootCmd.AddCommand(commands.StartTimerCmd)
	rootCmd.AddCommand(commands.StopTimerCmd)

	commands.StartTimerCmd.Flags().BoolVarP(&defFlag, "default", "d", false, "Use default project id.")
	// user.go
	rootCmd.AddCommand(commands.GetUserCmd)
	rootCmd.AddCommand(commands.AddPartTimeCmd)
	// report.go
	rootCmd.AddCommand(commands.CheckBalanceCmd)
	rootCmd.AddCommand(commands.VerfiyMonthCmd)

	// config.go
	rootCmd.AddCommand(commands.UpdateCmd)
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
