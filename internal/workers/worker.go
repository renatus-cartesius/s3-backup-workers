package workers

import (
	"archive/tar"
	"bytes"
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

	var buf bytes.Buffer
	for job := range b.Jobs {
		log.Println("Worker", b.Id, ": Doing backup", job.Src)

		err := packTarGz(job.Src, &buf)
		if err != nil {
			log.Fatalln(err)
			return err
		}

		fileToWrite, err := os.OpenFile(fmt.Sprint(job.Dst, filepath.Base(job.Src), ".tar.gz"), os.O_CREATE|os.O_RDWR, 0600)
		defer fileToWrite.Close()

		if err != nil {
			log.Fatalln(err)
			panic(err)
		}
		if _, err := io.Copy(fileToWrite, &buf); err != nil {
			log.Fatalln(err)
			panic(err)
		}
		buf.Reset()
	}

	log.Println("Worker", b.Id, ": died succesfully")
	return nil
}

func packTarGz(src string, buf io.Writer) error {
	tw := tar.NewWriter(buf)
	// zt := gzip.NewWriter(tw)

	var header *tar.Header

	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err = tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(path)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.Open(path)
			defer func() {
				err := data.Close()
				if err != nil {
					panic(err)
				}
			}()

			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error on filepath", src, ":", err)
	}

	if err := tw.Close(); err != nil {
		return err
	}

	// if err := zr.Close(); err != nil {
	// 	return err
	// }

	return nil
}
