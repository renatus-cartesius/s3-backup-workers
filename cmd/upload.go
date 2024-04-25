package cmd

import (
	"backup-workers/internal/workers"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to S3",
	Long:  `You can upload your files to S3 with compressing.`,
	Run: func(cmd *cobra.Command, args []string) {
		exec(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().String("src", "", "Path to dir to backup")
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

	log.Println("Reading config from", config)
	log.Println("Started uploading to S3")

	// Getting subdirs
	dirs, err := os.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan *workers.BackupJob, len(dirs))
	var wg sync.WaitGroup

	for w := 1; w <= 12; w++ {
		wg.Add(1)
		worker := workers.NewBackupWorker(w, jobs, &wg)
		go worker.StartWork()
	}

	go func() {
		for _, d := range dirs {
			if d.IsDir() {
				job := workers.NewBackupJob(src + "/" + d.Name())
				jobs <- job
			}
		}
		close(jobs)
	}()

	wg.Wait()
}
