package cmd

import (
	"aion/service"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rankCmd)
}

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Rank Player Info",
	Run: func(cmd *cobra.Command, args []string) {
		rankService := service.NewRankService()
		err := rankService.Run()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	},
}
