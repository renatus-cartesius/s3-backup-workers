package workers

import (
	"log"
	"math/rand"
	"time"
)

type BackupWorker struct {
}

func NewBackupWorker() *BackupWorker {
	return &BackupWorker{}
}

func (b *BackupWorker) Do(job *BackupJob, done chan int) error {
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	log.Println("Doing job ", job.Id, ": backuping and uploading", job.Path)
	done <- 1
	return nil
}
