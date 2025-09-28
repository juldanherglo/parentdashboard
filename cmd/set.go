package cmd

import (
	"fmt"
	"parentdashboard/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set times",
	Long:  `set times`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cliConfig := api.CliConfig{}

		if err := viper.Unmarshal(&cliConfig); err != nil {
			return err
		}

		return cliConfig.SetTimes()
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	setCmd.PersistentFlags().String("set-file-name", "", "Filename of the json file to set new time settings.")
	if err := viper.BindPFlag("set-file-name", setCmd.PersistentFlags().Lookup("set-file-name")); err != nil {
		fmt.Print(err)
	}

	setCmd.PersistentFlags().String("csrf-token", "", "CSRF_TOKEN for authentication of PUT request.")
	if err := viper.BindPFlag("csrf-token", setCmd.PersistentFlags().Lookup("csrf-token")); err != nil {
		fmt.Print(err)
	}
}
