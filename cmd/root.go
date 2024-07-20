package cmd

import (
	"aion/config"
	"aion/model"
	"aion/zlog"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "aion [command]",
	Short: "Aion Chatlog Analyse System",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "conf", "c", "app.ini", "config file")
}

func initConfig() {
	if configFile == "" {
		rootCmd.Help()
		return
	}
	err := config.Init(configFile)
	if err != nil {
		rootCmd.PrintErrln(err)
		return
	}
	zlog.Init()
	model.Init(false)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
