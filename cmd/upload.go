package cmd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"sync"

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
	Process() error
}

type prefixStage struct {
	in  chan []byte
	out chan []byte

	prefix string
}

func (ps *prefixStage) Process() error {
	defer close(ps.out)
	for data := range ps.in {

		modifiedData := fmt.Sprintf("%s{%s}", ps.prefix, data)

		ps.out <- []byte(modifiedData)

	}
	return nil
}

func newPrefixStage(wg *sync.WaitGroup, in chan []byte, prefix string) chan []byte {
	out := make(chan []byte)

	ps := &prefixStage{
		in:     in,
		out:    out,
		prefix: prefix,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ps.Process()
	}()

	return out
}

// type pipeline struct {
// 	wg        *sync.WaitGroup
// 	rootStage stage
// }

type gzipStage struct {
	in  chan []byte
	out chan []byte

	// buf *bytes.Buffer
}

func newGzipStage(wg *sync.WaitGroup, in chan []byte) chan []byte {
	out := make(chan []byte)

	// buf := bytes.NewBuffer(make([]byte, 0))
	// gw := gzip.NewWriter(buf)

	gs := &gzipStage{
		in:  in,
		out: out,
		// gw:  gw,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		gs.Process()
	}()

	return out
}

func (gs *gzipStage) Process() error {
	defer close(gs.out)
	for data := range gs.in {

		res := bytes.NewBuffer(data)

		gw := gzip.NewWriter(res)

		gw.Write(data)
		if err := gw.Flush(); err != nil {
			fmt.Println(err)
			return err
		}

		gs.out <- res.Bytes()

	}

	return nil
}

func execUpload(cmd *cobra.Command, args []string) {

	// var in chan []byte

	// prefixes := []string{"apple", "honey", "bread", "asdfasdf", "sdfasdfasdf234123"}
	// pipeIn := make(chan []byte)
	in := make(chan []byte)

	wg := &sync.WaitGroup{}

	// for i, p := range prefixes {

	// 	if i == 0 {
	// 		in = newPrefixStage(wg, pipeIn, p)
	// 		continue
	// 	}

	// 	in = newPrefixStage(wg, in, p)

	// }

	out1 := newGzipStage(wg, in)
	out := newPrefixStage(wg, out1, "SOME_LOOOONG_PREFIX")

	data, err := os.ReadFile("./lorem.test")
	if err != nil {
		panic(err)
	}

	in <- data

	fmt.Println(string(<-out))

}
