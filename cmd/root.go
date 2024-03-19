package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backup-workers",
	Short: "Simple backuping programm that uploads compressed files to S3 with telegram notifications",
	Long:  `This program developed for personal purposes of backuping.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "./backup.yml", "Path to config file")
}
