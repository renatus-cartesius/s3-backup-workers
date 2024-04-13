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
}

func exec(cmd *cobra.Command, args []string) {
	config, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Panic("error on loading config:", err)
	}
	log.Println("Reading config from", config)
	log.Println("Started uploading to S3")

	// Getting subdirs
	dirs, err := os.ReadDir("/")
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan *workers.BackupJob, len(dirs))
	var wg sync.WaitGroup

	for w := 1; w <= 20; w++ {
		wg.Add(1)
		worker := workers.NewBackupWorker(w, jobs, &wg)
		go worker.StartWork()
	}

	for i, d := range dirs {
		// if d.IsDir() {
		// 	worker := workers.NewBackupWorker()
		// 	job := workers.NewBackupJob(i, d.Name())
		// 	go worker.Do(job, done)
		// }
		job := workers.NewBackupJob(i, d.Name())
		jobs <- job
	}

	log.Println("Creating closer")
	log.Println("Closing jobs channel")
	wg.Wait()
	log.Println("Closing jobs channel")
	close(jobs)
	// go func() {
	// }()

	log.Println("Exiting main goroutine")
}
