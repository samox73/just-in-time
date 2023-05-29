package cmd

import (
	"os"

	"github.com/samox73/just-in-time/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var rootCmd = &cobra.Command{
	Use:   "jit",
	Short: "Just-In-Time",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		viper.SetEnvPrefix("just-in-time")
		viper.AutomaticEnv()

		_ = viper.ReadInConfig()

		slog.Info("This is Just-In-Time")
		return ui.RunTea()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
