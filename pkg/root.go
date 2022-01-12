package pkg

import (
	"fmt"
	"os"

	controller "github.com/Faagerholm/clockify-cli/pkg/Controller"
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
			controller.Menu()
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
	rootCmd.AddCommand(controller.AddKeyCmd)
	rootCmd.AddCommand(resetViperCmd)
	rootCmd.AddCommand(controller.SetupCmd)
	// menu.go
	rootCmd.AddCommand(controller.MenuCmd)
	// project.go
	rootCmd.AddCommand(controller.DefaultProjectCmd)
	rootCmd.AddCommand(controller.ListProjectsCmd)
	// entry.go
	rootCmd.AddCommand(controller.StartTimerCmd)
	rootCmd.AddCommand(controller.StopTimerCmd)

	controller.StartTimerCmd.Flags().BoolVarP(&defFlag, "default", "d", false, "Use default project id.")
	// user.go
	rootCmd.AddCommand(controller.GetUserCmd)
	rootCmd.AddCommand(controller.AddPartTimeCmd)
	// report.go
	rootCmd.AddCommand(controller.CheckBalanceCmd)
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
