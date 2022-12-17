package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile   string
	algorithm string

	rootCmd = &cobra.Command{
		Use:   "cobra-cli",
		Short: "A generator for Cobra based Applications",
		Long: `Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var cmdDemo = &cobra.Command{
	Use:   "start",
	Short: "Demo Post Msg",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aaaaa: %v\n", viper.GetString("algorithm"))
		fmt.Printf("str slice: %v\n", viper.GetStringSlice("backends"))
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&algorithm, "algorithm", "a", "default_algo", "load balancing algorithm to be used")
	viper.BindPFlag("algorithm", rootCmd.PersistentFlags().Lookup("algorithm"))

	addCommands()
}

func addCommands() {
	rootCmd.AddCommand(cmdDemo)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		fmt.Printf("HOME: %s\n", home)
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Failed reading config file:", err.Error())
	}
}
