package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "parentdashboard",
	Short: "Cli for alexa's parentdashboard settings",
	Long:  `Cli for alexa's parentdashboard settings`,
	// RunE: func(cmd *cobra.Command, args []string) error {
	// cliConfig := api.CliConfig{}

	// if err := viper.Unmarshal(&cliConfig); err != nil {
	// return err
	// }

	// return cliConfig.GetTimes()
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetTypeByDefaultValue(true)

	rootCmd.PersistentFlags().String("log-level", "info", "log-level can be: error, warn, info, debug")
	if err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		fmt.Print(err)
	}
}
