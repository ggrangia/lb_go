package cmd

import (
	"errors"
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
	"github.com/ggrangia/lb_go/pkg/lb_go/selection/wrr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "lb_go",
		Short: "A simple load balancer",
		Long:  "A simple load balancer",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start the loab balancer",
	RunE: func(cmd *cobra.Command, args []string) error {
		algo := viper.GetString("algorithm")
		port := viper.GetInt("port")
		hc := viper.GetInt("healthcheck")
		fmt.Printf("Selected algorithm: %v\n", viper.GetString("algorithm"))
		back_urls := viper.GetStringSlice("backends")
		backends := make([]*backend.Backend, len(back_urls))
		for i, b := range back_urls {
			backends[i] = backend.NewBackend(b)
		}
		weights := viper.GetIntSlice("weights")
		return start(backends, algo, time.Duration(hc), port, weights)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./lb_go.yaml)")
	rootCmd.PersistentFlags().StringP("algorithm", "a", "roundrobin", "load balancing algorithm to be used")
	rootCmd.PersistentFlags().IntP("healthcheck", "c", 5, "healthcheck timer")
	rootCmd.PersistentFlags().IntP("port", "p", 8080, "listening port")
	viper.BindPFlag("algorithm", rootCmd.PersistentFlags().Lookup("algorithm"))
	viper.BindPFlag("healthcheck", rootCmd.PersistentFlags().Lookup("healthcheck"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(cmdStart)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find the config in the running folder
		curr, err := os.Getwd()
		fmt.Printf("Current DIR: %s\n", curr)
		cobra.CheckErr(err)

		// Search config in current directory with name "lb_go"
		viper.AddConfigPath(curr)
		viper.SetConfigType("yaml")
		viper.SetConfigName("lb_go")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Failed reading config file:", err.Error())
		fmt.Println("****************")
		rootCmd.Help()
		os.Exit(0)
	}
}

var ErrMismatchLength = errors.New("mismatch length between backends and weights")

func start(backends []*backend.Backend, algo string, h time.Duration, port int, weights []int) error {
	var selector selection.Selector
	switch algo {
	case "roundrobin":
		selector = roundrobin.NewWithBackends(backends)
	case "randomselection":
		selector = randomselection.NewWithBackends(time.Now().UTC().UnixNano(), backends)
	case "wrr":
		w := wrr.New()
		lw := len(weights)
		lbs := len(backends)
		if lw == 0 {
			weights = make([]int, lbs)
		} else if lbs != lw {
			return ErrMismatchLength
		}
		for i, weight := range weights {
			w.AddWeightedBackend(backends[i], weight)
		}
		selector = w
	default:
		log.Fatalf("Unknown selection algorithm: %v", algo)
	}

	hc := healthcheck.New(selector, h)

	lb := lb_go.NewLb(selector, hc, port)
	lb.Start()
	return nil
}
