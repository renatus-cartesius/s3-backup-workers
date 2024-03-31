package workers

type BackupJob struct {
	Id   int
	Path string
}

func NewBackupJob(id int, path string) *BackupJob {
	return &BackupJob{
		Id:   id,
		Path: path,
	}
}
