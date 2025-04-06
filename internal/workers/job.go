package workers

type BackupJob struct {
	Src string
	Dst string
}

func NewBackupJob(src, dst string) *BackupJob {
	return &BackupJob{
		Src: src,
		Dst: dst,
	}
}
