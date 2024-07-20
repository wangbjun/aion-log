package cmd

import (
	"aion/model"
	"aion/service"
	"github.com/spf13/cobra"
)

var (
	logFile string
)

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.PersistentFlags().StringVarP(&logFile, "file", "f", "", "log file path")
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse Chatlog file",
	Run: func(cmd *cobra.Command, args []string) {
		if logFile == "" {
			cmd.Usage()
			return
		}
		model.Init(true)
		parser := service.NewParseService()
		err := parser.Run(logFile)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		service.NewClassifyService().Run()
		service.NewRankService().Run()
		service.NewTimelineService().Run()
	},
}
