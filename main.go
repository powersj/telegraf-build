package main

import (
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "telegraf-build",
	Short:            "Minifi telegraf build with only the required plugins",
	Long:             `Minifi telegraf build with only the required plugins`,
	Args:             args,
	PersistentPreRun: setup,
	RunE:             build,
	SilenceUsage:     true,
}

func init() {
	rootCmd.Flags().StringSliceVar(
		&buildConfigFiles, "config", []string{}, "config file to pick plugins from",
	)
	rootCmd.Flags().StringSliceVar(
		&buildAggregators, "aggregators", []string{}, "aggregator plugins to include",
	)
	rootCmd.Flags().StringSliceVar(
		&buildInputs, "inputs", []string{}, "input plugins to include",
	)
	rootCmd.Flags().StringSliceVar(
		&buildOutputs, "outputs", []string{}, "output plugins to include",
	)
	rootCmd.Flags().StringSliceVar(
		&buildProcessors, "processors", []string{}, "processor plugins to include",
	)

	rootCmd.Flags().StringVar(
		&sourceDirectory, "source", ".", "directory with Telegraf source and plugins directory",
	)
}

// ensure source directory exists and at least has the plugins directory.
func args(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(sourceDirectory); os.IsNotExist(err) {
		return errors.Wrap(err, "Telegraf source directory does not exist")
	}

	if _, err := os.Stat(path.Join(sourceDirectory, "plugins")); os.IsNotExist(err) {
		return errors.Wrap(err, "Telegraf plugins directory does not exist, please set --source")
	}

	return nil
}

// setup the logging format and level.
func setup(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
