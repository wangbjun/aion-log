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
		err := service.NewClassifyService().Run()
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}
		err = service.NewRankService().Run()
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}
		err = service.NewTimelineService().Run()
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}
	},
}
