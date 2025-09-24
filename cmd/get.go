package cmd

import (
	"parentdashboard/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get times",
	Long:  `get times`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cliConfig := api.CliConfig{}

		if err := viper.Unmarshal(&cliConfig); err != nil {
			return err
		}

		return cliConfig.GetTimes()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
