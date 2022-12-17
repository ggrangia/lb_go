package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ggrangia/lb_go/pkg/healthcheck"
	"github.com/ggrangia/lb_go/pkg/lb_go"
	"github.com/ggrangia/lb_go/pkg/lb_go/backend"
	"github.com/ggrangia/lb_go/pkg/lb_go/selection"
	"github.com/ggrangia/lb_go/pkg/lb_go/selection/randomselection"
	"github.com/ggrangia/lb_go/pkg/lb_go/selection/roundrobin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile   string
	algorithm string

	rootCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the loab balancer",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start the loab balancer",
	Run: func(cmd *cobra.Command, args []string) {
		algo := viper.GetString("algorithm")
		fmt.Printf("aaaaa: %v\n", viper.GetString("algorithm"))
		fmt.Printf("str slice: %v\n", viper.GetStringSlice("backends"))
		back_urls := viper.GetStringSlice("backends")
		backends := make([]*backend.Backend, len(back_urls))
		for i, b := range back_urls {
			backends[i] = backend.NewBackend(b)
		}
		start(backends, algo)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// FIXME: Change default config file name
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&algorithm, "algorithm", "a", "roundrobin", "load balancing algorithm to be used")
	rootCmd.PersistentFlags().IntP("healthcheck", "c", 5, "healthcheck timer")
	viper.BindPFlag("algorithm", rootCmd.PersistentFlags().Lookup("algorithm"))
	viper.BindPFlag("healthcheck", rootCmd.PersistentFlags().Lookup("healthcheck"))

	rootCmd.AddCommand(cmdStart)
}

func initConfig() {
	fmt.Println("AAAAAAAAAAAAAAA")
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

func start(backends []*backend.Backend, algo string) {
	var selector selection.Selector
	switch algo {
	case "roundrobin":
		selector = roundrobin.NewWithBackends(backends)
	case "randomselection":
		selector = randomselection.NewWithBackends(time.Now().UTC().UnixNano(), backends)
	default:
		log.Fatalf("Unknown selection algorithm: %v", algo)
	}

	hc := healthcheck.New(selector, 5)

	lb := lb_go.NewLb(selector, hc)
	lb.Start()
}
