package cmd

import (
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
}
