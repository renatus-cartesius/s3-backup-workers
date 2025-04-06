package cmd

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to S3",
	Long:  `You can upload your files to S3 with compressing.`,
	Run: func(cmd *cobra.Command, args []string) {
		execUpload(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

type stage interface {
	GetInput() io.WriteCloser
	CloseInput() error

	GetOutput() io.ReadCloser
	CloseOutput() error

	Next() stage
}

type prefixWriter struct {
	prefix string
	dstW   io.Writer
}

func (pw *prefixWriter) Write(p []byte) (int, error) {
	data := fmt.Sprintf("%s{%s}", pw.prefix, p)

	return pw.dstW.Write([]byte(data))
}

type prefixStage struct {
	out  io.ReadCloser
	in   *prefixWriter
	next stage
}

func (ps *prefixStage) GetInput() io.WriteCloser {
	panic("not implemented") // TODO: Implement
}

func (ps *prefixStage) CloseInput() error {
	panic("not implemented") // TODO: Implement
}

func (ps *prefixStage) GetOutput() io.ReadCloser {
	panic("not implemented") // TODO: Implement
}

func (ps *prefixStage) CloseOutput() error {
	panic("not implemented") // TODO: Implement
}

func (ps *prefixStage) Next() stage {
	panic("not implemented") // TODO: Implement
}

func newPrefixStage(prefix string, next stage) *prefixStage {
	r, w := io.Pipe()

	pw := &prefixWriter{
		prefix: prefix,
		dstW:   w,
	}

	return &prefixStage{
		out: r,
		in:  pw,
	}
}

type gzipStage struct {
	out  io.ReadCloser
	in   *gzip.Writer
	next stage
}

func (gs *gzipStage) Next() stage {
	return gs.next
}

func (gs *gzipStage) GetInput() io.WriteCloser {
	return gs.in
}

func (gs *gzipStage) GetOutput() io.ReadCloser {
	return gs.out
}

func (gs *gzipStage) CloseInput() error {
	return gs.in.Close()
}

func (gs *gzipStage) CloseOutput() error {
	return gs.out.Close()
}

func newGzipStage(next stage) *gzipStage {
	out, in := io.Pipe()

	w := gzip.NewWriter(in)

	return &gzipStage{
		in:   w,
		out:  out,
		next: next,
	}
}

type fileWriter struct {
	file   *os.File
	bypass io.WriteCloser
}

func (fs *fileStage) GetInput() io.WriteCloser {
	return fs.in
}

func (fs *fileStage) GetOutput() io.ReadCloser {
	return fs.out
}

func (fs *fileStage) CloseInput() error {
	return fs.in.Close()
}

func (fs *fileStage) CloseOutput() error {
	return fs.out.Close()
}

func (fw *fileWriter) Write(p []byte) (int, error) {
	_, err := fw.file.Write(p)
	if err != nil {
		return 0, err
	}
	return fw.bypass.Write(p)
}

func (fw *fileWriter) Close() error {
	if err := fw.bypass.Close(); err != nil {
		return err
	}
	return fw.file.Close()
}

type fileStage struct {
	out  io.ReadCloser
	in   *fileWriter
	next stage
}

func (fs *fileStage) Next() stage {
	return fs.next
}

func newFileStage(next stage) *fileStage {

	out, in := io.Pipe()

	file, err := os.Create(fmt.Sprintf("./testing-%s", uuid.NewString()))
	if err != nil {
		panic(err)
	}

	w := &fileWriter{
		file:   file,
		bypass: in,
	}

	return &fileStage{
		out:  out,
		in:   w,
		next: next,
	}
}

type stringReaderCustom struct {
	data string
}

func (src *stringReaderCustom) Read(dst []byte) (int, error) {
	fmt.Println("[+] reading from custom reader")
	return strings.NewReader(src.data).Read(dst)
}

// Close method needs for stringsReaderCustom implement io.ReadCloser
func (src *stringReaderCustom) Close() error {
	return nil
}

type stringStage struct {
	out  *stringReaderCustom
	next stage
}

func (ss *stringStage) Next() stage {
	return ss.next
}

func (ss *stringStage) GetOutput() io.ReadCloser {
	return ss.out
}

func (ss *stringStage) GetInput() io.WriteCloser {
	return nil
}

func (ss *stringStage) CloseInput() error {
	return nil
}

func (ss *stringStage) CloseOutput() error {
	return ss.out.Close()
}

func newStringStage(next stage) *stringStage {
	out := &stringReaderCustom{
		data: "Note that a single db_type override configuration applies to either nullable or non-nullable columns, but not both. If you want the same Go type to override in both cases, you’ll need to configure two overrides.",
	}

	return &stringStage{
		out:  out,
		next: next,
	}
}

type pipeline struct {
	wg        *sync.WaitGroup
	rootStage stage
}

func (p *pipeline) Process(source io.Reader) error {

	p.wg.Add(1)
	go func() {
		// defer p.rootStage.CloseInput()
		// defer p.rootStage.CloseOutput()
		defer p.wg.Done()
		io.Copy(p.rootStage.GetInput(), source) // Need no closing source
		p.rootStage.CloseInput()
		// p.rootStage.CloseOutput()
	}()

	end := proccessStage(p.rootStage, p.wg)

	p.wg.Add(1)
	go func() {
		// defer end.CloseInput()
		// defer end.CloseOutput()
		defer p.wg.Done()
		io.Copy(os.Stdout, end.GetOutput())
	}()

	p.wg.Wait()
	return nil
}

// proccessStage processing pipeline stages going through linked list and return last stage
func proccessStage(s stage, wg *sync.WaitGroup) stage {
	fmt.Println("Proccessing", reflect.TypeOf(s))

	if s.Next() != nil {

		fmt.Println("Found non-end stage")

		wg.Add(1)
		go func() {
			// defer s.CloseInput()
			// defer s.CloseOutput()
			defer wg.Done()

			fmt.Println(reflect.TypeOf(s), "->", reflect.TypeOf(s.Next()))
			if n, err := io.Copy(s.Next().GetInput(), s.GetOutput()); err != nil {
				fmt.Printf("\n%s %s %s", n, err, reflect.TypeOf(s))
			}
			// s.CloseInput()
			s.Next().CloseInput()
			s.CloseOutput()

			fmt.Println("\nDone on", reflect.TypeOf(s))
		}()

		return proccessStage(s.Next(), wg)
	}

	// Returning last stage (assuing last stage has `nil` next field)
	return s
}

func execUpload(cmd *cobra.Command, args []string) {

	input := strings.NewReader("Note that a single db_type override configuration applies to either nullable or non-nullable columns, but not both. If you want the same Go type to override in both cases, you’ll need to configure two overrides.")

	sp := newPrefixStage("sdfsdf", nil)

	go func() {
		defer sp.out.Close()
		sp.in.Write([]byte("some data that will be prefixed"))
	}()
	io.Copy(os.Stdout, sp.out)

	// Starting from last
	sf := newFileStage(nil)
	// sgz1 := newGzipStage(sf)
	sgz := newGzipStage(sf)
	// ss := newStringStage(sgz)

	// pipeline := []stage{ss, sgz, sf}

	pipeline := &pipeline{
		wg:        &sync.WaitGroup{},
		rootStage: sgz,
	}

	pipeline.Process(input)

	// 	wg.Add(1)
	// 	go func() {
	// 		defer pipeline[i-1].CloseOutput()
	// 		defer pipeline[i-1].CloseInput()
	// 		defer pipeline[i].CloseInput()
	// 		defer pipeline[i].CloseOutput()
	// 		defer wg.Done()
	// 		fmt.Println("Copy from ", i-1, "to", i)
	// 		io.Copy(pipeline[i].GetInput(), pipeline[i-1].GetOutput())
	// 		fmt.Println("Done from ", i-1, "to", i)
	// 	}()
	// }

	// wg := sync.WaitGroup{}

	// // Pipeline start - writing input to rootStage
	// wg.Add(1)
	// go func() {
	// 	defer sgz.out.Close()
	// 	defer sgz.in.Close()
	// 	defer wg.Done()
	// 	io.Copy(sgz.in, input)
	// 	fmt.Println("executed 3")
	// }()

	// // Pipeline stage processing - chaining inner stages outputs and inputs
	// wg.Add(1)
	// go func() {
	// 	defer sgz1.in.Close()
	// 	defer sgz1.out.Close()
	// 	defer wg.Done()
	// 	io.Copy(sgz1.in, sgz.out)
	// 	fmt.Println("executed 2")
	// }()
	// wg.Add(1)
	// go func() {
	// 	defer sf.in.Close()
	// 	defer sf.out.Close()
	// 	defer wg.Done()
	// 	io.Copy(sf.in, sgz1.out)
	// 	fmt.Println("executed 2")
	// }()

	// // Pipeline end - writing last stage output to somewhere
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	io.Copy(os.Stdout, sf.out)
	// 	fmt.Println("executed 1")
	// }()

	// wg.Wait()
}
