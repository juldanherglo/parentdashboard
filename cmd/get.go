package cmd

import (
	"fmt"
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

	getCmd.PersistentFlags().String("get-file-name", "", "Filename of the json file to write current time gettings into.")
	if err := viper.BindPFlag("get-file-name", getCmd.PersistentFlags().Lookup("get-file-name")); err != nil {
		fmt.Print(err)
	}
}
