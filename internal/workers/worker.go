package workers

import (
	"log"
	"sync"
	"time"
)

type BackupWorker struct {
	Id        int
	Jobs      <-chan *BackupJob
	WaitGroup *sync.WaitGroup
}

func NewBackupWorker(id int, jobs <-chan *BackupJob, wg *sync.WaitGroup) *BackupWorker {
	return &BackupWorker{
		Id:        id,
		Jobs:      jobs,
		WaitGroup: wg,
	}
}

func (b *BackupWorker) StartWork() error {
	defer b.WaitGroup.Done()

	for job := range b.Jobs {
		time.Sleep(time.Second * 1)
		log.Println("Worker", b.Id, ": Doing job ", job.Id, ": backuping and uploading", job.Path)
	}

	log.Println("Worker", b.Id, ": died succesfully")
	return nil
}
