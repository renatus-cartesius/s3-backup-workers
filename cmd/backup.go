/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"backup-workers/internal/workers"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup target dir, separatly by subdirs",
	Long:  `Concurrently packing and compressing subdirs to target dir`,
	Run: func(cmd *cobra.Command, args []string) {
		exec(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().String("src", "", "Path to dir to backup")
	backupCmd.Flags().String("dst", "", "Path to destination dir, where backup will store")
}

func exec(cmd *cobra.Command, args []string) {
	config, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Panic("error on loading config:", err)
	}

	src, err := cmd.Flags().GetString("src")
	if err != nil {
		log.Panic("error reading src:", err)
	}

	dst, err := cmd.Flags().GetString("dst")
	if err != nil {
		log.Panic("error reading dst:", err)
	}

	log.Println("Reading config from", config)
	log.Println("Started uploading to S3")

	// Getting subdirs
	dirs, err := os.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan *workers.BackupJob, len(dirs))
	var wg sync.WaitGroup

	for w := 1; w <= runtime.NumCPU(); w++ {
		wg.Add(1)
		worker := workers.NewBackupWorker(w, jobs, &wg)
		go worker.StartWork()
	}

	go func() {
		for _, d := range dirs {
			if d.IsDir() {
				job := workers.NewBackupJob(src+"/"+d.Name(), dst)
				jobs <- job
			}
		}
		close(jobs)
	}()

	wg.Wait()
}
