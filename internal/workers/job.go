package workers

type BackupJob struct {
	Path string
}

func NewBackupJob(path string) *BackupJob {
	return &BackupJob{
		Path: path,
	}
}
