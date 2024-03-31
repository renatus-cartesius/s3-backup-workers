package workers

import (
	"log"
	"time"
)

type BackupWorker struct {
	Id   int
	Jobs <-chan BackupJob
	Done chan<- struct{}
}

func NewBackupWorker(id int, jobs <-chan BackupJob, done chan<- struct{}) *BackupWorker {
	return &BackupWorker{
		Id:   id,
		Jobs: jobs,
		Done: done,
	}
}

func (b *BackupWorker) Do() error {
	for job := range b.Jobs {
		time.Sleep(time.Second * 1)
		log.Println("Worker", b.Id, ": Doing job ", job.Id, ": backuping and uploading", job.Path)
	}
	return nil
}
