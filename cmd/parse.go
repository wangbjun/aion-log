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
		err := model.DB().Exec("truncate table aion_player_chat_log").Error
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		err = model.DB().Exec("truncate table aion_player_info").Error
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		parser := service.NewParseService()
		err = parser.Run(logFile)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		service.NewClassifyService().Run()
		service.NewRankService().Run()
	},
}
