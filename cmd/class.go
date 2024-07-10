package cmd

import (
	"aion/service"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(classCmd)
}

var classCmd = &cobra.Command{
	Use:   "class",
	Short: "Classify Player Info",
	Run: func(cmd *cobra.Command, args []string) {
		service.NewClassifyService().Run()
		service.NewRankService().Run()
		service.NewTimelineService().Run()
	},
}
