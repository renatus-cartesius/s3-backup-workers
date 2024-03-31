package cmd

import (
	"backup-workers/internal/workers"
	"log"
	"os"

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

	done := make(chan struct{})
	jobs := make(chan workers.BackupJob, len(dirs))

	for w := 1; w <= 20; w++ {
		worker := workers.NewBackupWorker(w, jobs, done)
		go worker.Do()
	}

	for i, d := range dirs {
		// if d.IsDir() {
		// 	worker := workers.NewBackupWorker()
		// 	job := workers.NewBackupJob(i, d.Name())
		// 	go worker.Do(job, done)
		// }
		job := workers.NewBackupJob(i, d.Name())
		jobs <- *job
	}

	for range dirs {
		<-done
	}

}
