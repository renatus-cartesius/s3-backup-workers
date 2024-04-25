package workers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
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
		log.Println("Worker", b.Id, ": Doing backup", job.Path)

		var buf bytes.Buffer
		err := packTarGz(job.Path, &buf)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		fileToWrite, err := os.OpenFile(fmt.Sprint("/backup/", filepath.Base(job.Path), ".tar.gz"), os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			log.Fatalln(err)
			panic(err)
		}
		if _, err := io.Copy(fileToWrite, &buf); err != nil {
			log.Fatalln(err)
			panic(err)
		}
	}

	log.Println("Worker", b.Id, ": died succesfully")
	return nil
}

func packTarGz(src string, buf io.Writer) error {
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(path)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
			data.Close()
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error on filepath", src, ":", err)
	}

	if err := tw.Close(); err != nil {
		return err
	}

	if err := zr.Close(); err != nil {
		return err
	}

	return nil
}
