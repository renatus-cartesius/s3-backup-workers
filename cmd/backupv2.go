/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// backupv2Cmd represents the backupv2 command
var backupv2Cmd = &cobra.Command{
	Use:   "backupv2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: execBackupv2,
}

type mockTarProcessor struct{}

// type mockEncryptStage struct{}
// type mockFileStage struct{}
// type mockS3Stage struct{}

type dataPipelineStage interface {
	// SetInput(io.)
	Execute(io.Writer) io.Reader
}

type mockPipeline struct {
	stages []dataPipelineStage

	// gzip    mockGzipStage
	// tar     mockTarStage
	// encrypt mockEncryptStage
}

func (mp *mockPipeline) Process() {
	// for _, stage := range mp.stages {

	// }
}

// func (mpf *mockPipelineFacade)

type mockCore struct {
	pipeline mockPipeline
}

func init() {

	rootCmd.AddCommand(backupv2Cmd)
	// for _, cs := range []int{30, 2, 10} {

	// chunk := make([]byte, 2<<10)

	// var buf bytes.Buffer

	// fmt.Println("Writing to buffer...")

	// start := time.Now()

	// io.CopyBuffer(&buf, file, chunk)

	// fmt.Printf("Ended writing, collapsed: %v\n", time.Since(start))
	// fmt.Println()
	// }

	// n, err := r.Read(buf)
	// if err != nil {
	// 	panic(err)
	// }
	// buf = buf[:n]

	// for n != 0 {
	// 	fmt.Println(string(buf))
	// 	n, err := r.Read(buf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		panic(err)
	// 	}
	// 	buf = buf[:n]
	// }

	// gw := gzip.NewWriter(os.Stdout)

	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// 	time.Sleep(time.Second * 2)
	// }

	// gw.Write([]byte("sdfsdf\x00"))

	// fmt.Println(buf.String())

}

type mockGzipStage struct {
	in  io.Reader
	out io.Writer
	gz  *gzip.Writer
}

func NewMockGzipStage(in io.Reader, out io.Writer) *mockGzipStage {
	return &mockGzipStage{
		in:  in,
		out: out,
		gz:  gzip.NewWriter(out),
	}
}

func (mgs *mockGzipStage) Process() {
	_, err := io.Copy(mgs.gz, mgs.in)
	defer mgs.gz.Flush()
	if err != nil {
		panic(err)
	}
}

// type simpleStage struct {
// 	writer io.Writer
// }

// func (s *simpleStage) Resolve(r io.Reader) (out io.Reader) {

type writer struct {
}

func (w *writer) Write(p []byte) (n int, err error) {
	return copy(p, []byte("test")), nil
}

// func Resolve(r io.Reader) (out io.Reader) {
// wr :=
// _, err := io.Copy()
// }

// func (mgs *mockGzipStage) Read(p []byte){
// 	_, err := io.Copy(mgs.out, mgs.in)
// 	defer mgs.out.Flush()
// 	if err != nil {
// 		panic(err)
// 	}
// }

func execBackupv2(cmd *cobra.Command, args []string) {

	// file, err := os.Open("/home/renatus/tmp_files/random2.txt")
	// if err != nil {
	// 	panic(err)
	// }

	w := &writer{}
	str := make([]byte, 100)
	w.Write(str)

	fmt.Println(string(str))

	mgs1 := NewMockGzipStage(
		strings.NewReader("ssdf8as9f 8as9f8u sajf ;k2q341-02i 4-0i a-isdf asdfdfsfsdf\n"),
		os.Stdout,
	)

	// mgs2 := NewMockGzipStage(
	// 	mgs1
	// )

	mgs1.Process()

	r := io.MultiReader(
		strings.NewReader("HEAD\n"),
		strings.NewReader("BODY\n"),
	)

	io.Copy(os.Stdout, r)

	// str := strings.NewReader("sfsdf")
	// gw := gzip.NewWriter(os.Stdin)
	// io.Copy(gw, str)
	// gw.Flush()
}
